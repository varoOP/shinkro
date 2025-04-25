package server

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/varoOP/shinkro/internal/mapping"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/domain"
)

type Server struct {
	log            zerolog.Logger
	config         *domain.Config
	animeService   anime.Service
	mappingService mapping.Service
}

func NewServer(log zerolog.Logger, config *domain.Config, animeSvc anime.Service, mappingSvc mapping.Service) *Server {
	return &Server{
		log:            log.With().Str("module", "server").Logger(),
		config:         config,
		animeService:   animeSvc,
		mappingService: mappingSvc,
	}
}

func (s *Server) Start() error {
	err := s.animeService.UpdateAnime()
	if err != nil {
		return err
	}

	if _, err := s.mappingService.Get(context.Background()); errors.Is(err, sql.ErrNoRows) {
		_ = s.mappingService.Store(context.Background(), &domain.MapSettings{
			TVDBEnabled:       false,
			TMDBEnabled:       false,
			CustomMapTMDBPath: "",
			CustomMapTVDBPath: "",
		})
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
