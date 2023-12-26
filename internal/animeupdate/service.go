package animeupdate

import (
	"context"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
)

type Service interface {
	Store(ctx context.Context, animeupdate *domain.AnimeUpdate) error
	GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error)
	UpdateAnimeList(ctx context.Context, anime *domain.AnimeUpdate, event domain.PlexEvent) error
}

type service struct {
	log            zerolog.Logger
	repo           domain.AnimeUpdateRepo
	animeService   anime.Service
	mapService     mapping.Service
	malauthService malauth.Service
}

func NewService(log zerolog.Logger, repo domain.AnimeUpdateRepo, animeSvc anime.Service, mapSvc mapping.Service, malauthSvc malauth.Service) Service {
	return &service{
		log:            log.With().Str("module", "animeUpdate").Logger(),
		repo:           repo,
		animeService:   animeSvc,
		mapService:     mapSvc,
		malauthService: malauthSvc,
	}
}

func (s *service) Store(ctx context.Context, animeupdate *domain.AnimeUpdate) error {
	return s.repo.Store(ctx, animeupdate)
}

func (s *service) GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error) {
	return s.repo.GetByID(ctx, req)
}

func (s *service) UpdateAnimeList(ctx context.Context, anime *domain.AnimeUpdate, event domain.PlexEvent) error {
    switch event {
    case domain.PlexRateEvent:
        return s.handleRateEvent(ctx, anime)
    case domain.PlexScrobbleEvent:
        return s.handleScrobbleEvent(ctx, anime)
    }
    return nil
}

func (s *service) handleRateEvent(ctx context.Context, anime *domain.AnimeUpdate) error {
    if anime.SourceDB == domain.MAL {
        anime.MALId = anime.SourceId
        return s.updateAndStore(ctx, anime, anime.UpdateRating)
    }

    animeMap, err := s.mapService.CheckForAnimeinMap(ctx, anime)
    if err == nil {
        anime.MALId = animeMap.Malid
        return s.updateAndStore(ctx, anime, anime.UpdateRating)
    }

    if anime.SeasonNum == 1 {
        return s.updateFromDBAndStore(ctx, anime, anime.UpdateRating)
    }

    return err
}

func (s *service) handleScrobbleEvent(ctx context.Context, anime *domain.AnimeUpdate) error {
    if anime.SourceDB == domain.MAL {
        anime.MALId = anime.SourceId
        return s.updateAndStore(ctx, anime, anime.UpdateWatchStatus)
    }

    animeMap, err := s.mapService.CheckForAnimeinMap(ctx, anime)
    if err == nil {
        anime.MALId = animeMap.Malid
        anime.EpisodeNum = animeMap.CalculateEpNum(anime.EpisodeNum)
        return s.updateAndStore(ctx, anime, anime.UpdateWatchStatus)
    }

    if anime.SeasonNum == 1 {
        return s.updateFromDBAndStore(ctx, anime, anime.UpdateWatchStatus)
    }

    return err
}

func (s *service) updateAndStore(ctx context.Context, anime *domain.AnimeUpdate, updateFunc func(context.Context, *mal.Client) error) error {
    client, err := s.malauthService.GetMalClient(ctx)
    if err != nil {
        return err
    }

    if err := updateFunc(ctx, client); err != nil {
        return err
    }

    return s.Store(ctx, anime)
}

func (s *service) updateFromDBAndStore(ctx context.Context, anime *domain.AnimeUpdate, updateFunc func(context.Context, *mal.Client) error) error {
    req := &domain.GetAnimeRequest{
        IDtype: anime.SourceDB,
        Id:     anime.SourceId,
    }

    animeFromDB, err := s.animeService.GetByID(ctx, req)
    if err != nil {
        return err
    }

    anime.MALId = animeFromDB.MALId
    return s.updateAndStore(ctx, anime, updateFunc)
}

