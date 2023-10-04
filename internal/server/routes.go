package server

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
			r.Use(Auth(cfg))
			r.Route("/plex", func(r chi.Router) {
				r.Use(OnlyAllowPost, middleware.AllowContentType("application/json", "multipart/form-data"), ParsePlexPayload, CheckPlexPayload(cfg))
				r.Post("/", Plex(db, cfg, &log, n))
			})
		})

		r.Route("/malauth", func(r chi.Router) {
			r.Use(BasicAuth(cfg.Username, cfg.Password))
			r.With(CheckMalAuth(db)).Get("/", MalAuth(cfg))
			r.Post("/login", MalAuthLogin())
			r.Get("/callback", MalAuthCallback(cfg, db, &log))
			r.Get("/status", MalAuthStatus(cfg, db))
		})

		r.NotFound(NotFound(cfg))
	})

	return r
}
