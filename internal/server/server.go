package server

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/domain"
)

type Server struct {
	log          zerolog.Logger
	config       *domain.Config
	animeService anime.Service
}

func NewServer(log zerolog.Logger, config *domain.Config, animeSvc anime.Service) *Server {
	return &Server{
		log:          log.With().Str("module", "server").Logger(),
		config:       config,
		animeService: animeSvc,
	}
}

func (s *Server) Start() error {
	err := s.animeService.UpdateAnime()
	if err != nil {
		return err
	}

	c := cron.New(cron.WithLocation(time.UTC))
	_, err = c.AddFunc("0 1 * * MON", func() {
		s.animeService.UpdateAnime()
	})

	if err != nil {
		return err
	}

	c.Start()
	return nil
}
