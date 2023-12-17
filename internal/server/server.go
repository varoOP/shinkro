package server

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/plex"
)

type Server struct {
	log        zerolog.Logger
	config     *domain.Config
	plexService plex.Service

	stopWG sync.WaitGroup
	lock   sync.Mutex
}

func NewServer(log zerolog.Logger, config *domain.Config, plexSvc plex.Service) *Server {
	return &Server{
		log: log.With().Str("module", "server").Logger(),
		config: config,
		plexService: plexSvc,
	}
}
