package server

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/mapping"
)

type Server struct {
	log          zerolog.Logger
	config       *domain.Config
	animeService anime.Service
	mapService   mapping.Service
}

func NewServer(log zerolog.Logger, config *domain.Config, animeSvc anime.Service, mapSvc mapping.Service) *Server {
	return &Server{
		log:          log.With().Str("module", "server").Logger(),
		config:       config,
		animeService: animeSvc,
		mapService:   mapSvc,
	}
}

func (s *Server) Start() error {
	err := s.animeService.UpdateAnime()
	if err != nil {
		return err
	}

	err, mapLoaded := s.mapService.CheckLocalMaps()
	if err != nil {
		s.log.Fatal().Err(err).Msg("Unable to load local custom mapping.")
	}

	if mapLoaded {
		s.log.Info().Msg("Loaded local custom mapping successfully.")
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
