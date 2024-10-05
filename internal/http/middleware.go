package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/domain"
)

// func auth(cfg *domain.Config) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			log := hlog.FromRequest(r)
// 			if !isAuthorized(cfg.ApiKey, r.URL.Query()) && !isAuthorized(cfg.ApiKey, r.Header) {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				log.Error().Err(errors.New("ApiKey invalid")).Msg("")
// 				log.Debug().Str("query", fmt.Sprintf("%v", r.URL.Query())).Str("headers", fmt.Sprintf("%v", r.Header)).Msg("")
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		}

// 		return http.HandlerFunc(fn)
// 	}
// }

func (s Server) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := r.Header.Get("Shinkro-Api-Key"); token != "" {
			// check header
			if !s.apiService.ValidateAPIKey(r.Context(), token) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

		} else if key := r.URL.Query().Get("apiKey"); key != "" {
			// check query param like ?apiKey=TOKEN
			if !s.apiService.ValidateAPIKey(r.Context(), key) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		} else {
			// check session
			session, err := s.cookieStore.Get(r, "user_session")
			if err != nil {
				s.log.Error().Err(err).Msgf("could not get session from cookieStore")
				session.Values["authenticated"] = false

				// MaxAge<0 means delete cookie immediately
				session.Options.MaxAge = -1
				session.Options.Path = s.config.BaseUrl

				if err := session.Save(r, w); err != nil {
					s.log.Error().Err(err).Msgf("could not store session: %s", r.RemoteAddr)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			// Check if user is authenticated
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				s.log.Warn().Msg("session not authenticated")

				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if created, ok := session.Values["created"].(int64); ok {
				// created is a unix timestamp MaxAge is in seconds
				maxAge := time.Duration(session.Options.MaxAge) * time.Second
				expires := time.Unix(created, 0).Add(maxAge)

				if time.Until(expires) <= 7*24*time.Hour { // 7 days
					s.log.Info().Msgf("Cookie is expiring in less than 7 days on %s - extending session", expires.Format("2006-01-02 15:04:05"))

					session.Values["created"] = time.Now().Unix()

					// Call session.Save as needed - since it writes a header (the Set-Cookie
					// header), making sure you call it before writing out a body is important.
					// https://github.com/gorilla/sessions/issues/178#issuecomment-447674812
					if err := session.Save(r, w); err != nil {
						s.log.Error().Err(err).Msgf("could not store session: %s", r.RemoteAddr)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}

			ctx := context.WithValue(r.Context(), sessionkey, session)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// func basicAuth(username, password string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			user, pass, ok := r.BasicAuth()
// 			if !ok || user != username || pass != password {
// 				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
// 				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
// 				return
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// func checkMalAuth(db *database.DB) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			log := hlog.FromRequest(r)
// 			// client, _ := malauth.NewOauth2Client(r.Context(), db)
// 			c := mal.NewClient(&http.Client{})
// 			_, _, err := c.User.MyInfo(r.Context())
// 			if err == nil {
// 				w.Write([]byte("Authentication with myanimelist is successful."))
// 				log.Trace().Msg("user already authenticated")
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		}

// 		return http.HandlerFunc(fn)
// 	}
// }

// func onlyAllowPost(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		log := hlog.FromRequest(r)
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			log.Error().Err(errors.New("method not allowed")).Msg("")
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	}

// 	return http.HandlerFunc(fn)
// }

func parsePlexPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := hlog.FromRequest(r)

		sourceType := contentType(r)
		log.Trace().Str("sourceType", string(sourceType)).Msg("")

		payload, err := parsePayloadBySourceType(w, r, sourceType)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), domain.PlexPayload, payload)
		log.Debug().Str("parsedPlexPayload", fmt.Sprintf("%+v", payload)).Msg("")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func checkPlexPayload(cfg *domain.Config) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			log := hlog.FromRequest(r)
// 			p := r.Context().Value(domain.PlexPayload).(*plex.PlexWebhook)
// 			aa := "Accepted"
// 			if !isPlexUser(p, cfg) {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				log.Error().Err(errors.New("unauthorized plex user")).
// 					Str("plexUserReceived", p.Account.Title).
// 					Str("AuthorizedPlexUser", cfg.PlexUser).
// 					Msg("")

// 				return
// 			}

// 			if !isEvent(p) {
// 				http.Error(w, aa, http.StatusAccepted)
// 				log.Trace().Err(errors.New("incorrect event")).
// 					Str("event", p.Event).
// 					Str("allowedEvents", "media.scrobble, media.rate").
// 					Msg("")

// 				return
// 			}

// 			if !isAnimeLibrary(p, cfg) {
// 				http.Error(w, aa, http.StatusAccepted)
// 				log.Error().Err(errors.New("not an anime library")).
// 					Str("library received", p.Metadata.LibrarySectionTitle).
// 					Str("anime libraries", strings.Join(cfg.AnimeLibraries, ",")).
// 					Msg("")

// 				return
// 			}

// 			allowed, agent := isMetadataAgent(p)
// 			if !allowed {
// 				http.Error(w, aa, http.StatusAccepted)
// 				log.Debug().Err(errors.New("unsupported metadata agent")).
// 					Str("guid", string(p.Metadata.GUID.GUID)).
// 					Str("supported metadata agents", "HAMA, MyAnimeList.bundle, Plex Series, Plex Movie").
// 					Msg("")

// 				return
// 			}

// 			mediaTypeOk := mediaType(p)
// 			if !mediaTypeOk {
// 				http.Error(w, aa, http.StatusAccepted)
// 				log.Debug().Err(errors.New("unsupported media type")).
// 					Str("media type", p.Metadata.Type).
// 					Str("supported media types", "episode, movie").
// 					Msg("")

// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), domain.Agent, agent)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		}

// 		return http.HandlerFunc(fn)
// 	}
// }
