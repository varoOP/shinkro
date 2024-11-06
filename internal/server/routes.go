package server

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
)

func NewRouter(cfg *domain.Config, db *database.DB, n *domain.Notification, log zerolog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(hlog.NewHandler(log))
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Int("status", status).
			Dur("duration", duration).
			Msg("Request processed")
	}))

	baseUrl, err := url.JoinPath("/", cfg.BaseUrl)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	r.Route(baseUrl, func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Use(auth(cfg))
			r.Route("/plex", func(r chi.Router) {
				r.Use(onlyAllowPost, middleware.AllowContentType("application/json", "multipart/form-data"), parsePlexPayload, checkPlexPayload(cfg))
				r.Post("/", plexHandler(db, cfg, &log, n))
			})
		})

		r.Route("/malauth", func(r chi.Router) {
			r.Use(basicAuth(cfg.Username, cfg.Password))
			r.With(checkMalAuth(db)).Get("/", malAuth(cfg))
			r.Post("/login", malAuthLogin())
			r.Get("/callback", malAuthCallback(cfg, db, &log))
			r.Get("/status", malAuthStatus(cfg, db))
		})

		r.NotFound(notFound(cfg))
	})

	return r
}
