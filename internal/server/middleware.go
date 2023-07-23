package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
		log := hlog.FromRequest(r)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "recevied bad request", http.StatusBadRequest)
			log.Trace().Msg("received bad request")
			return
		}

		ps := r.PostFormValue("payload")
		log.Trace().RawJSON("rawPlexPayload", []byte(ps)).Msg("")
		payload, err := plex.NewPlexWebhook([]byte(ps))
		if err != nil {
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
			if !isPlexUser(p, cfg) {
				log.Debug().Err(errors.New("unauthorized plex user")).
					Str("plexUserReceived", p.Account.Title).
					Str("AuthorizedPlexUser", cfg.PlexUser).
					Msg("")

				return
			}

			if !isEvent(p) {
				log.Trace().Err(errors.New("incorrect event")).
					Str("event", p.Event).
					Str("allowedEvents", "media.scrobble, media.rate").
					Msg("")

				return
			}

			if !isAnimeLibrary(p, cfg) {
				log.Debug().Err(errors.New("not an anime library")).
					Str("library received", p.Metadata.LibrarySectionTitle).
					Str("anime libraries", strings.Join(cfg.AnimeLibraries, ",")).
					Msg("")

				return
			}

			if !isMetadataAgent(p) {
				log.Debug().Err(errors.New("unsupported metadata agent")).
					Str("guid", string(p.Metadata.GUID.GUID)).
					Str("supported metadata agents", "HAMA, MyAnimeList.bundle").
					Msg("")

				return
			}

			mediaTypeOk, title := mediaType(p)
			if !mediaTypeOk {
				log.Debug().Err(errors.New("unsupported media type")).
					Str("media type", p.Metadata.Type).
					Str("supported media types", "episode, movie").
					Msg("")

				return
			}

			ctx := context.WithValue(r.Context(), domain.MediaTitle, title)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
