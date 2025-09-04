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
	Store(ctx context.Context, userID int, animeupdate *domain.AnimeUpdate) error
	GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error)
	UpdateAnimeList(ctx context.Context, userID int, anime *domain.AnimeUpdate, event domain.PlexEvent) error
	Count(ctx context.Context) (int, error)
	GetRecentUnique(ctx context.Context, userID int, limit int) ([]*domain.AnimeUpdate, error)
	GetByPlexID(ctx context.Context, plexID int64) (*domain.AnimeUpdate, error)
	GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]*domain.AnimeUpdate, error)
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

func (s *service) Store(ctx context.Context, userID int, animeupdate *domain.AnimeUpdate) error {
	return s.repo.Store(ctx, userID, animeupdate)
}

func (s *service) GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error) {
	return s.repo.GetByID(ctx, req)
}

func (s *service) UpdateAnimeList(ctx context.Context, userID int, anime *domain.AnimeUpdate, event domain.PlexEvent) error {
	switch event {
	case domain.PlexRateEvent:
		return s.handleRateEvent(ctx, userID, anime)
	case domain.PlexScrobbleEvent:
		return s.handleScrobbleEvent(ctx, userID, anime)
	}
	return nil
}

func (s *service) handleRateEvent(ctx context.Context, userID int, anime *domain.AnimeUpdate) error {
	return s.handleEvent(ctx, userID, anime, anime.UpdateRating, false)
}

func (s *service) handleScrobbleEvent(ctx context.Context, userID int, anime *domain.AnimeUpdate) error {
	return s.handleEvent(ctx, userID, anime, anime.UpdateWatchStatus, true)
}

func (s *service) handleEvent(ctx context.Context, userID int, anime *domain.AnimeUpdate, updateFunc func(context.Context, *mal.Client) error, isScrobble bool) error {
	if anime.SourceDB == domain.MAL {
		anime.MALId = anime.SourceId
		return s.updateAndStore(ctx, userID, anime, updateFunc)
	}

	convertedAnime := s.convertAniDBToTVDB(ctx, anime)
	animeMap, err := s.mapService.CheckForAnimeinMap(ctx, convertedAnime)
	if err == nil {
		anime.MALId = animeMap.Malid
		if isScrobble {
			anime.EpisodeNum = animeMap.CalculateEpNum(anime.EpisodeNum)
		}
		return s.updateAndStore(ctx, userID, anime, updateFunc)
	}

	if anime.SeasonNum == 1 {
		return s.updateFromDBAndStore(ctx, userID, anime, updateFunc)
	}

	return err
}

func (s *service) updateAndStore(ctx context.Context, userID int, anime *domain.AnimeUpdate, updateFunc func(context.Context, *mal.Client) error) error {
	client, err := s.malauthService.GetMalClient(ctx, userID)
	if err != nil {
		return err
	}

	if err := updateFunc(ctx, client); err != nil {
		return err
	}
	s.log.Info().Interface("status", anime.ListStatus).Msg("MyAnimeList Updated Successfully")
	return s.Store(ctx, userID, anime)
}

func (s *service) updateFromDBAndStore(ctx context.Context, userID int, anime *domain.AnimeUpdate, updateFunc func(context.Context, *mal.Client) error) error {
	req := &domain.GetAnimeRequest{
		IDtype: anime.SourceDB,
		Id:     anime.SourceId,
	}

	animeFromDB, err := s.animeService.GetByID(ctx, req)
	if err != nil {
		return err
	}

	anime.MALId = animeFromDB.MALId
	return s.updateAndStore(ctx, userID, anime, updateFunc)
}

func (s *service) convertAniDBToTVDB(ctx context.Context, anime *domain.AnimeUpdate) *domain.AnimeUpdate {
	if anime.SourceDB != domain.AniDB {
		return anime
	}

	req := &domain.GetAnimeRequest{
		IDtype: anime.SourceDB,
		Id:     anime.SourceId,
	}

	aa, err := s.animeService.GetByID(ctx, req)
	if err != nil {
		return anime
	}

	newAnime := *anime
	if aa.TVDBId > 0 {
		newAnime.SourceDB = domain.TVDB
		newAnime.SourceId = aa.TVDBId
		s.log.Debug().Int("converted tvdbId", aa.TVDBId).Msg("Converted Anime to TVDB")
	} else {
		return anime
	}

	return &newAnime
}

func (s *service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

func (s *service) GetRecentUnique(ctx context.Context, userID int, limit int) ([]*domain.AnimeUpdate, error) {
	return s.repo.GetRecentUnique(ctx, userID, limit)
}

func (s *service) GetByPlexID(ctx context.Context, plexID int64) (*domain.AnimeUpdate, error) {
	return s.repo.GetByPlexID(ctx, plexID)
}

func (s *service) GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]*domain.AnimeUpdate, error) {
	if len(plexIDs) == 0 {
		return []*domain.AnimeUpdate{}, nil
	}
	return s.repo.GetByPlexIDs(ctx, plexIDs)
}
