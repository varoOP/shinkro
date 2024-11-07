package http

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
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
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(hlog.NewHandler(s.log))
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Int("status", status).
			Dur("duration", duration).
			Msg("Request processed")
	}))

	baseUrll, err := url.JoinPath("/", s.config.BaseUrl)
	if err != nil {
		s.log.Error().Err(err).Msg("")
	}

	baseUrl := s.config.BaseUrl
	if !strings.HasPrefix(baseUrl, "/") {
		baseUrl = "/" + baseUrl
	}
	if baseUrl != "/" && !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	c := cors.New(cors.Options{
		AllowCredentials:   true,
		AllowedMethods:     []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowOriginFunc:    func(origin string) bool { return true },
		OptionsPassthrough: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	r.Use(c.Handler)

	encoder := encoder{}

	apiRouter := chi.NewRouter()
	apiRouter.Route("/auth", newAuthHandler(encoder, s.log, s, s.config, s.cookieStore, s.authService).Routes)

	apiRouter.Group(func(r chi.Router) {
		r.Use(s.IsAuthenticated)
		r.Route("/plex", newPlexHandler(encoder, s.plexService).Routes)
		r.Route("/malauth", newmalauthHandler(encoder, s.malauthService).Routes)
		r.Route("/keys", newAPIKeyHandler(encoder, s.apiService).Routes)
	})

	// Mount API routes under baseUrl + "api"
	r.Mount(baseUrl+"api", apiRouter)

	// Create a separate web router for the SPA and static files
	webRouter := chi.NewRouter()
	web.RegisterHandler(webRouter, s.version, baseUrl)

	// Mount the web router under the baseUrl
	r.Mount(baseUrll, webRouter)

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
