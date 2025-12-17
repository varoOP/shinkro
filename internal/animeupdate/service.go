package animeupdate

import (
	"context"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rs/zerolog"
	"github.com/pkg/errors"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
)

type Service interface {
	Store(ctx context.Context, animeupdate *domain.AnimeUpdate) error
	GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error)
	UpdateAnimeList(ctx context.Context, anime *domain.AnimeUpdate, event domain.PlexEvent) error
	Count(ctx context.Context) (int, error)
	GetRecentUnique(ctx context.Context, limit int) ([]*domain.AnimeUpdate, error)
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
	return s.handleEvent(ctx, anime, false)
}

func (s *service) handleScrobbleEvent(ctx context.Context, anime *domain.AnimeUpdate) error {
	return s.handleEvent(ctx, anime, true)
}

func (s *service) handleEvent(ctx context.Context, anime *domain.AnimeUpdate, isScrobble bool) error {
	if anime.SourceDB == domain.MAL {
		anime.MALId = anime.SourceId
		return s.updateAndStore(ctx, anime, isScrobble)
	}

	convertedAnime := s.convertAniDBToTVDB(ctx, anime)
	animeMap, err := s.mapService.CheckForAnimeinMap(ctx, convertedAnime)
	if err == nil {
		anime.MALId = animeMap.Malid
		if isScrobble {
			anime.EpisodeNum = animeMap.CalculateEpNum(anime.EpisodeNum)
		}
		return s.updateAndStore(ctx, anime, isScrobble)
	}

	if anime.SeasonNum == 1 {
		return s.updateFromDBAndStore(ctx, anime, isScrobble)
	}

	return err
}

func (s *service) updateAndStore(ctx context.Context, anime *domain.AnimeUpdate, isScrobble bool) error {
	client, err := s.malauthService.GetMalClient(ctx)
	if err != nil {
		return err
	}

	// Fetch current anime list details from MAL API
	if err := s.fetchAnimeDetails(ctx, client, anime); err != nil {
		return err
	}

	// Update MAL based on event type
	if isScrobble {
		if err := s.updateWatchStatus(ctx, client, anime); err != nil {
			return err
		}
	} else {
		if err := s.updateRating(ctx, client, anime); err != nil {
			return err
		}
	}

	s.log.Info().Interface("status", anime.ListStatus).Msg("MyAnimeList Updated Successfully")
	return s.Store(ctx, anime)
}

func (s *service) updateFromDBAndStore(ctx context.Context, anime *domain.AnimeUpdate, isScrobble bool) error {
	req := &domain.GetAnimeRequest{
		IDtype: anime.SourceDB,
		Id:     anime.SourceId,
	}

	animeFromDB, err := s.animeService.GetByID(ctx, req)
	if err != nil {
		return err
	}

	s.log.Debug().Int("malId", animeFromDB.MALId).Msg("Anime from DB")
	if animeFromDB.MALId == 0 {
		return errors.New("Anime not found in database, update mapping")
	}

	anime.MALId = animeFromDB.MALId
	return s.updateAndStore(ctx, anime, isScrobble)
}

// fetchAnimeDetails calls MAL API to get current anime list details
func (s *service) fetchAnimeDetails(ctx context.Context, client *mal.Client, anime *domain.AnimeUpdate) error {
	aa, _, err := client.Anime.Details(ctx, anime.MALId, mal.Fields{"num_episodes", "title", "main_picture{medium,large}", "my_list_status{status,num_times_rewatched,num_episodes_watched}"})
	if err != nil {
		return err
	}

	details := domain.BuildListDetailsFromMALResponse(
		aa.MyListStatus.Status,
		aa.MyListStatus.NumTimesRewatched,
		aa.NumEpisodes,
		aa.MyListStatus.NumEpisodesWatched,
		aa.Title,
		aa.MainPicture.Medium,
	)
	anime.UpdateListDetails(details)

	return nil
}

// updateRating calls MAL API to update rating and updates domain with result
func (s *service) updateRating(ctx context.Context, client *mal.Client, anime *domain.AnimeUpdate) error {
	l, _, err := client.Anime.UpdateMyListStatus(ctx, anime.MALId, mal.Score(anime.Plex.Rating))
	if err != nil {
		return err
	}

	anime.UpdateRatingWithStatus(*l)
	return nil
}

// updateWatchStatus calls MAL API to update watch status and updates domain with result
func (s *service) updateWatchStatus(ctx context.Context, client *mal.Client, anime *domain.AnimeUpdate) error {
	options, err := anime.BuildWatchStatusOptions()
	if err != nil {
		return err
	}

	l, _, err := client.Anime.UpdateMyListStatus(ctx, anime.MALId, options...)
	if err != nil {
		return err
	}

	anime.UpdateWatchStatusWithStatus(*l)
	return nil
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

func (s *service) GetRecentUnique(ctx context.Context, limit int) ([]*domain.AnimeUpdate, error) {
	return s.repo.GetRecentUnique(ctx, limit)
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
