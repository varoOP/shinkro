package http

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/web"
)

type Server struct {
	log            zerolog.Logger
	db             *database.DB
	config         *domain.Config
	cookieStore    *sessions.CookieStore
	version        string
	commit         string
	date           string
	plexService    plexService
	malauthService malauthService
	apiService     apikeyService
	authService    authService
}

func NewServer(log zerolog.Logger, config *domain.Config, db *database.DB, version string, commit string, date string, plexSvc plexService, malauthSvc malauthService, apiSvc apikeyService, authSvc authService) Server {
	return Server{
		log:     log.With().Str("module", "http").Logger(),
		config:  config,
		db:      db,
		version: version,
		commit:  commit,
		date:    date,

		plexService:    plexSvc,
		malauthService: malauthSvc,
		apiService:     apiSvc,
		authService:    authSvc,
	}
}

func (s Server) Open() error {
	addr := fmt.Sprintf("%v:%v", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := http.Server{
		Handler:           s.Handler(),
		ReadHeaderTimeout: time.Second * 15,
	}

	s.log.Info().Msgf("Starting server. Listening on %s", listener.Addr().String())

	return server.Serve(listener)
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(hlog.NewHandler(s.log))
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Int("status", status).
			Dur("duration", duration).
			Msg("Request processed")
	}))

	baseUrl, err := url.JoinPath("/", s.config.BaseUrl)
	if err != nil {
		s.log.Error().Err(err).Msg("")
	}

	encoder := encoder{}

	r.Route(baseUrl, func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Route("/auth", newAuthHandler(encoder, s.log, s, s.config, s.cookieStore, s.authService).Routes)

			r.Group(func(r chi.Router) {
				r.Use(s.IsAuthenticated)
				r.Route("/plex", newPlexHandler(encoder, s.plexService).Routes)
				r.Route("/malauth", newmalauthHandler(encoder, s.malauthService).Routes)
				r.Route("/keys", newAPIKeyHandler(encoder, s.apiService).Routes)
			})
		})
		// r.Use(basicAuth(s.config.Username, s.config.Password))
		// r.With(checkMalAuth(s.db)).Get("/", malAuth(s.config))
		// r.Post("/login", malAuthLogin())
		// r.Get("/callback", malAuthCallback(s.config, s.db, &s.log))
		// r.Get("/status", malAuthStatus(s.config, s.db))

		// r.NotFound(notFound(s.config))
	})

	// serve the web
	web.RegisterHandler(r, s.version, s.config.BaseUrl)
	return r
}

// func (s Server) index(w http.ResponseWriter, r *http.Request) {
// 	p := web.IndexParams{
// 		Title:   "Dashboard",
// 		Version: s.version,
// 		BaseUrl: s.config.BaseUrl,
// 	}
// 	web.Index(w, p)
// }
