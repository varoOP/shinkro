package animeupdate

import (
	"context"

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
		if anime.SourceDB == domain.MAL {
			anime.MALId = anime.SourceId
			client, err := s.malauthService.GetMalClient(ctx)
			if err != nil {
				return err
			}

			err = anime.UpdateRating(ctx, client)
			if err != nil {
				return err
			}

			err = s.Store(ctx, anime)
			if err != nil {
				return err
			}

			return nil
		}

		animeMap, err := s.mapService.CheckForAnimeinMap(ctx, anime)
		if err == nil {
			anime.MALId = animeMap.Malid
			client, err := s.malauthService.GetMalClient(ctx)
			if err != nil {
				return err
			}
			err = anime.UpdateRating(ctx, client)
			if err != nil {
				return err
			}

			err = s.Store(ctx, anime)
			if err != nil {
				return err
			}

			return nil
		}

		if anime.SeasonNum == 1 {
			req := &domain.GetAnimeRequest{
				IDtype: anime.SourceDB,
				Id:     anime.SourceId,
			}

			animeFromDB, err := s.animeService.GetByID(ctx, req)
			if err != nil {
				return err
			}

			anime.MALId = animeFromDB.MALId
			client, err := s.malauthService.GetMalClient(ctx)
			if err != nil {
				return err
			}

			err = anime.UpdateRating(ctx, client)
			if err != nil {
				return err
			}

			err = s.Store(ctx, anime)
			if err != nil {
				return err
			}

			return nil
		}

	case domain.PlexScrobbleEvent:
		if anime.SourceDB == domain.MAL {
			anime.MALId = anime.SourceId
			client, err := s.malauthService.GetMalClient(ctx)
			if err != nil {
				return err
			}

			done, err := anime.UpdateWatchStatus(ctx, client)
			if err != nil {
				return err
			}

			if done {
				err = s.Store(ctx, anime)
				if err != nil {
					return err
				}
			}

			return nil
		}

		animeMap, err := s.mapService.CheckForAnimeinMap(ctx, anime)
		if err == nil {
			anime.MALId = animeMap.Malid
			anime.EpisodeNum = animeMap.CalculateEpNum(anime.EpisodeNum)
			client, err := s.malauthService.GetMalClient(ctx)
			if err != nil {
				return err
			}
			done, err := anime.UpdateWatchStatus(ctx, client)
			if err != nil {
				return err
			}

			if done {
				err = s.Store(ctx, anime)
				if err != nil {
					return err
				}
			}

			return nil
		}

		if anime.SeasonNum == 1 {
			req := &domain.GetAnimeRequest{
				IDtype: anime.SourceDB,
				Id:     anime.SourceId,
			}

			animeFromDB, err := s.animeService.GetByID(ctx, req)
			if err != nil {
				return err
			}

			anime.MALId = animeFromDB.MALId
			client, err := s.malauthService.GetMalClient(ctx)
			if err != nil {
				return err
			}

			done, err := anime.UpdateWatchStatus(ctx, client)
			if err != nil {
				return err
			}

			if done {
				err = s.Store(ctx, anime)
				if err != nil {
					return err
				}
			}

			return nil
		}

	}

	return nil
}
