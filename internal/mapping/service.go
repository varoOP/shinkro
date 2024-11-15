package mapping

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	NewMap(ctx context.Context) (*domain.AnimeMap, error)
	CheckForAnimeinMap(ctx context.Context, anime *domain.AnimeUpdate) (*domain.AnimeMapDetails, error)
}

type service struct {
	log  zerolog.Logger
	repo domain.MappingRepo
}

func NewService(log zerolog.Logger, repo domain.MappingRepo) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) NewMap(ctx context.Context) (*domain.AnimeMap, error) {
	mapping := &domain.Mapping{}

	mapSettings, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	localTVDB, localTMDB := mapSettings.LocalMapsExist()

	if localTVDB || localTMDB {
		err = mapping.LoadLocalMaps(localTVDB, localTMDB)
		if err != nil {
			return nil, err
		}
	}

	if !localTVDB || !localTMDB {
		err = mapping.LoadCommunityMaps(ctx, localTVDB, localTMDB)
		if err != nil {
			return nil, err
		}
	}

	return &mapping.AnimeMap, nil
}

func (s *service) CheckForAnimeinMap(ctx context.Context, anime *domain.AnimeUpdate) (*domain.AnimeMapDetails, error) {
	animeMap, err := s.NewMap(ctx)
	if err != nil {
		return nil, err
	}

	switch anime.Plex.Metadata.Type {
	case domain.PlexMovie:
		inMap, animeMovie := animeMap.AnimeMovies.CheckMap(anime.SourceId)
		if inMap {
			return &domain.AnimeMapDetails{
				Malid: animeMovie.MALID,
				Start: 0,
			}, nil
		}
	case domain.PlexEpisode:
		inMap, animeTV := animeMap.AnimeTVShows.CheckMap(anime.SourceId, anime.SeasonNum, anime.EpisodeNum)
		if inMap {
			return &domain.AnimeMapDetails{
				Malid:      animeTV.Malid,
				Start:      animeTV.Start,
				UseMapping: animeTV.UseMapping,
			}, nil
		}
	}

	return nil, errors.New("anime not found in map")
}
