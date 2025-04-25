package plex

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/animeupdate"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
	"github.com/varoOP/shinkro/internal/plexsettings"
)

type Service interface {
	Store(ctx context.Context, plex *domain.Plex) error
	// FindAll(ctx context.Context) ([]*domain.Plex, error)
	Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error)
	ProcessPlex(ctx context.Context, plex *domain.Plex) error
	// ProcessPlexScrobbleEvent(plex *domain.Plex) error
	// Delete(ctx context.Context, req *domain.DeletePlexRequest) error
}

type service struct {
	log                zerolog.Logger
	repo               domain.PlexRepo
	plexettingsService plexsettings.Service
	animeService       anime.Service
	mapService         mapping.Service
	malauthService     malauth.Service
	animeUpdateService animeupdate.Service
}

func NewService(log zerolog.Logger, plexsettingsSvc plexsettings.Service, repo domain.PlexRepo, animeSvc anime.Service, mapSvc mapping.Service, malauthSvc malauth.Service, animeUpdateSvc animeupdate.Service) Service {
	return &service{
		log:                log.With().Str("module", "plex").Logger(),
		repo:               repo,
		plexettingsService: plexsettingsSvc,
		animeService:       animeSvc,
		mapService:         mapSvc,
		malauthService:     malauthSvc,
		animeUpdateService: animeUpdateSvc,
	}
}

func (s *service) Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error) {
	return s.repo.Get(ctx, req)
}

func (s *service) Store(ctx context.Context, plex *domain.Plex) error {
	return s.repo.Store(ctx, plex)
}

func (s *service) ProcessPlex(ctx context.Context, plex *domain.Plex) error {
	a, err := s.extractSourceIdForAnime(ctx, plex)
	if err != nil {
		return err
	}

	err = s.animeUpdateService.UpdateAnimeList(ctx, a, plex.Event)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) extractSourceIdForAnime(ctx context.Context, plex *domain.Plex) (*domain.AnimeUpdate, error) {
	plexSettings, err := s.plexettingsService.Get(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return nil, err
	}

	agent, err := plex.CheckPlex(plexSettings)
	if err != nil {
		s.log.Debug().Err(err).Msg("")
		return nil, err
	}

	source, id, err := s.getSourceIDFromAgent(ctx, plex, agent)
	if err != nil {
		return nil, err
	}

	a := plex.SetAnimeFields(source, id)
	return &a, nil
}

func (s *service) getSourceIDFromAgent(ctx context.Context, p *domain.Plex, agent domain.PlexSupportedAgents) (domain.PlexSupportedDBs, int, error) {
	switch agent {
	case domain.HAMA, domain.MALAgent:
		return p.Metadata.GUID.HamaMALAgent(agent)
	case domain.PlexAgent:
		return s.plexettingsService.HandlePlexAgent(ctx, p)
	}
	return "", 0, nil
}
