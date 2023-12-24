package plex

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/animeupdate"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
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
	config             *domain.Config
	animeService       anime.Service
	mapService         mapping.Service
	malauthService     malauth.Service
	animeUpdateService animeupdate.Service
}

func NewService(log zerolog.Logger, config *domain.Config, repo domain.PlexRepo, animeSvc anime.Service, mapSvc mapping.Service, malauthSvc malauth.Service, animeUpdateSvc animeupdate.Service) Service {
	return &service{
		log:                log.With().Str("module", "plex").Logger(),
		repo:               repo,
		config:             config,
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
	anime, err := s.extractSourceIdForAnime(ctx, plex)
	if err != nil {
		return err
	}

	err = s.animeUpdateService.UpdateAnimeList(ctx, anime, plex.Event)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) extractSourceIdForAnime(ctx context.Context, plex *domain.Plex) (*domain.AnimeUpdate, error) {
	agent, err := plex.CheckPlex(s.config)
	if err != nil {
		s.log.Debug().Err(err).Msg("")
		return nil, err
	}

	source, id, err := plex.GetSourceIDFromAgent(agent, s.config)
	if err != nil {
		return nil, err
	}

	if source == domain.AniDB && plex.Metadata.ParentIndex > 1 {
		req := &domain.GetAnimeRequest{
			IDtype: source,
			Id:     id,
		}

		a, err := s.animeService.GetByID(ctx, req)
		if err != nil {
			return nil, err
		}

		source = domain.TVDB
		id = a.TVDBId
	}

	anime := plex.SetAnimeFields(source, id)
	return &anime, nil
}
