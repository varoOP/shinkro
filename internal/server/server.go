package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/domain"
)

type Server struct {
	config *domain.Config
	notify *domain.Notification
	db     *database.DB
	log    zerolog.Logger
}

func NewServer(cfg *domain.Config, n *domain.Notification, db *database.DB, log *zerolog.Logger) *Server {
	return &Server{
		config: cfg,
		notify: n,
		db:     db,
		log:    log.With().Str("module", "server").Logger(),
	}
}

func (s *Server) Start() {
	router := NewRouter(s.config, s.db, s.notify, s.log)
	addr := fmt.Sprintf("%v:%v", s.config.Host, s.config.Port)
	s.log.Info().Msgf("starting http server on %v", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		s.log.Fatal().Err(err).Msg("failed to start http server")
	}
}
