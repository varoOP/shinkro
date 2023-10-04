package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/pkg/errors"

	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/tautulli"
	"github.com/varoOP/shinkro/pkg/plex"
)

func Auth(cfg *domain.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := hlog.FromRequest(r)
			if !isAuthorized(cfg.ApiKey, r.URL.Query()) && !isAuthorized(cfg.ApiKey, r.Header) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Error().Err(errors.New("ApiKey invalid")).Msg("")
				log.Debug().Str("query", fmt.Sprintf("%v", r.URL.Query())).Str("headers", fmt.Sprintf("%v", r.Header)).Msg("")
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func BasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || user != username || pass != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CheckMalAuth(db *database.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := hlog.FromRequest(r)
			client, _ := malauth.NewOauth2Client(r.Context(), db)
			c := mal.NewClient(client)
			_, _, err := c.User.MyInfo(r.Context())
			if err == nil {
				w.Write([]byte("Authentication with myanimelist is successful."))
				log.Trace().Msg("user already authenticated")
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func OnlyAllowPost(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := hlog.FromRequest(r)
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Error().Err(errors.New("method not allowed")).Msg("")
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ParsePlexPayload(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := hlog.FromRequest(r)
		var (
			ps      string
			payload *plex.PlexWebhook
		)

		sourceType := contentType(r)
		switch sourceType {
		case "plexWebhook":
			err := r.ParseMultipartForm(0)
			if err != nil {
				http.Error(w, "recevied bad request", http.StatusBadRequest)
				log.Trace().Err(err).Msg("received bad request")
				return
			}

			ps = r.PostFormValue("payload")
			log.Trace().Str("sourceType", sourceType).RawJSON("rawPlexPayload", []byte(ps)).Msg("")
			payload, err = plex.NewPlexWebhook([]byte(ps))
			if err != nil {
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				log.Error().Err(err).Msg("unable to unmarshal plex payload")
				return
			}

		case "tautulli":
			ps, err := readRequest(r)
			if err != nil {
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				log.Trace().Err(err).Msg(InternalServerError)
				return
			}

			log.Trace().Str("sourceType", sourceType).RawJSON("rawPlexPayload", []byte(ps)).Msg("")
			payload, err = tautulli.ToPlex([]byte(ps))
			if err != nil {
				http.Error(w, InternalServerError, http.StatusInternalServerError)
				log.Error().Err(err).Msg("unable to convert tautulli to plex webhook")
				return
			}

		default:
			log.Error().Str("sourceType", sourceType).Msg("sourceType not supported")
			return
		}

		ctx := context.WithValue(r.Context(), domain.PlexPayload, payload)
		log.Debug().Str("parsedPlexPayload", fmt.Sprintf("%+v", payload)).Msg("")
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func CheckPlexPayload(cfg *domain.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := hlog.FromRequest(r)
			p := r.Context().Value(domain.PlexPayload).(*plex.PlexWebhook)
			br := "bad request"
			if !isPlexUser(p, cfg) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Error().Err(errors.New("unauthorized plex user")).
					Str("plexUserReceived", p.Account.Title).
					Str("AuthorizedPlexUser", cfg.PlexUser).
					Msg("")

				return
			}

			if !isEvent(p) {
				http.Error(w, br, http.StatusBadRequest)
				log.Trace().Err(errors.New("incorrect event")).
					Str("event", p.Event).
					Str("allowedEvents", "media.scrobble, media.rate").
					Msg("")

				return
			}

			if !isAnimeLibrary(p, cfg) {
				http.Error(w, br, http.StatusBadRequest)
				log.Error().Err(errors.New("not an anime library")).
					Str("library received", p.Metadata.LibrarySectionTitle).
					Str("anime libraries", strings.Join(cfg.AnimeLibraries, ",")).
					Msg("")

				return
			}

			allowed, agent := isMetadataAgent(p)
			if !allowed {
				http.Error(w, br, http.StatusBadRequest)
				log.Debug().Err(errors.New("unsupported metadata agent")).
					Str("guid", string(p.Metadata.GUID.GUID)).
					Str("supported metadata agents", "HAMA, MyAnimeList.bundle, Plex Series, Plex Movie").
					Msg("")

				return
			}

			mediaTypeOk := mediaType(p)
			if !mediaTypeOk {
				http.Error(w, br, http.StatusBadRequest)
				log.Debug().Err(errors.New("unsupported media type")).
					Str("media type", p.Metadata.Type).
					Str("supported media types", "episode, movie").
					Msg("")

				return
			}

			ctx := context.WithValue(r.Context(), domain.Agent, agent)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
