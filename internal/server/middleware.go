package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/domain"
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

func ParsePlexPayload(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ps string
		log := hlog.FromRequest(r)
		sourceType, err := contentType(r)
		if err != nil {
			http.Error(w, "Unsupported content type", http.StatusNotAcceptable)
			log.Trace().Err(err).Str("Content-Type", sourceType).Msg("received unsupported content type")
			return
		}

		switch sourceType {
		case "plexWebhook":
			err = r.ParseMultipartForm(32 << 20)
			if err != nil {
				http.Error(w, "recevied bad request", http.StatusBadRequest)
				log.Trace().Err(err).Msg("received bad request")
				return
			}

			ps = r.PostFormValue("payload")

		case "tautulli":
			ps, err = readRequest(r)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Trace().Err(err).Msg("internal server error")
				return
			}
		}

		log.Trace().RawJSON("rawPlexPayload", []byte(ps)).Msg("")
		payload, err := plex.NewPlexWebhook([]byte(ps))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Error().Err(err).Msg("unable to unmarshal plex payload")
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
				log.Debug().Err(errors.New("unauthorized plex user")).
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
				log.Debug().Err(errors.New("not an anime library")).
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
