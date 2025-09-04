package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/domain"
)

func (s Server) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userContext context.Context = r.Context()
		
		if token := r.Header.Get("Shinkro-Api-Key"); token != "" {
			// check header
			if !s.apiService.ValidateAPIKey(r.Context(), token) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			
			// Get user ID from API key and add to context
			userID, err := s.apiService.GetUserIDByAPIKey(r.Context(), token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			userContext = context.WithValue(r.Context(), "api_user_id", userID)

		} else if key := r.URL.Query().Get("apiKey"); key != "" {
			// check query param like ?apiKey=TOKEN
			if !s.apiService.ValidateAPIKey(r.Context(), key) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			
			// Get user ID from API key and add to context
			userID, err := s.apiService.GetUserIDByAPIKey(r.Context(), key)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			userContext = context.WithValue(r.Context(), "api_user_id", userID)
		} else {
			// check session
			session, err := s.cookieStore.Get(r, "user_session")
			if err != nil {
				s.log.Error().Err(err).Msgf("could not get session from cookieStore")
				session.Values["authenticated"] = false

				// MaxAge<0 means delete cookie immediately
				session.Options.MaxAge = -1
				session.Options.Path = s.config.Config.BaseUrl

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

			userContext = context.WithValue(userContext, sessionkey, session)
		}

		next.ServeHTTP(w, r.WithContext(userContext))
	})
}

func parsePlexPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := hlog.FromRequest(r)

		sourceType := contentType(r)
		msg := fmt.Sprintf("sourceType: %s", string(sourceType))
		log.Debug().Msg(msg)

		payload, err := parsePayloadBySourceType(w, r, sourceType)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), domain.PlexPayload, payload)
		log.Debug().Str("parsedPlexPayload", fmt.Sprintf("%+v", payload)).Msg("")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s Server) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session from context (set by IsAuthenticated middleware)
		session, ok := r.Context().Value(sessionkey).(*sessions.Session)
		if !ok || session == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Get username from session
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Check if user is admin
		user, err := s.authService.FindByUsername(r.Context(), username)
		if err != nil {
			s.log.Error().Err(err).Msgf("RequireAdmin: could not find user: %s", username)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !user.Admin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
