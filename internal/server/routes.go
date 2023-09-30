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
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("You shall not pass!"))
		})

		r.Route("/plex", func(r chi.Router) {
			r.Use(OnlyAllowPost)
			r.Use(middleware.AllowContentType("application/json", "multipart/form-data"))
			r.Use(Auth(cfg))
			r.Use(ParsePlexPayload)
			r.Use(CheckPlexPayload(cfg))
			r.Post("/", Plex(db, cfg, &log, n))
		})
	})

	return r
}
