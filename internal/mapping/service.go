package mapping

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/varoOP/shinkro/internal/domain"
	"gopkg.in/yaml.v3"
)

type Service interface {
	NewAnimeMaps(ctx context.Context) (*domain.AnimeMap, error)
	CheckLocalMaps() (error, bool)
	CheckForAnimeinMap(ctx context.Context, anime *domain.AnimeUpdate) (*domain.AnimeMapDetails, error)
}

type service struct {
	log    zerolog.Logger
	config *domain.Config
}

func NewService(log zerolog.Logger, config *domain.Config) Service {
	return &service{
		log:    log,
		config: config,
	}
}

func (s *service) NewAnimeMaps(ctx context.Context) (*domain.AnimeMap, error) {
	s.config.LocalMapsExist()
	err := loadCommunityMaps(ctx, s.config)
	if err != nil {
		return nil, err
	}

	err = loadLocalMaps(s.config)
	if err != nil {
		return nil, err
	}

	return &domain.AnimeMap{
		AnimeTVShows: s.config.TVDBMalMap,
		AnimeMovies:  s.config.TMDBMalMap,
	}, nil
}

func (s *service) CheckLocalMaps() (error, bool) {
	loadLocalMaps(s.config)
	localMapLoaded := false
	if s.config.CustomMapTVDB {
		if err := validateYaml(domain.TVDBSchema, s.config.TVDBMalMap); err != nil {
			return err, false
		}

		localMapLoaded = true
	}

	if s.config.CustomMapTMDB {
		if err := validateYaml(domain.TMDBSchema, s.config.TMDBMalMap); err != nil {
			return err, false
		}

		localMapLoaded = true
	}

	return nil, localMapLoaded
}

func (s *service) CheckForAnimeinMap(ctx context.Context, anime *domain.AnimeUpdate) (*domain.AnimeMapDetails, error) {
	animeMap, err := s.NewAnimeMaps(ctx)
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

func loadCommunityMaps(ctx context.Context, cfg *domain.Config) error {
	if !cfg.CustomMapTVDB {
		s := &domain.AnimeTVShows{}
		respTVDB, err := domain.GetWithContext(ctx, string(domain.CommunityMapTVDB))
		if err != nil {
			return err
		}

		err = readYamlHTTP(respTVDB, s)
		if err != nil {
			return err
		}

		cfg.TVDBMalMap = s
	}

	if !cfg.CustomMapTMDB {
		am := &domain.AnimeMovies{}
		respTMDB, err := domain.GetWithContext(ctx, string(domain.CommunityMapTMDB))
		if err != nil {
			return err
		}

		err = readYamlHTTP(respTMDB, am)
		if err != nil {
			return err
		}

		cfg.TMDBMalMap = am
	}

	return nil
}

func loadLocalMaps(cfg *domain.Config) error {
	if cfg.CustomMapTVDB {
		s := &domain.AnimeTVShows{}
		fTVDB, err := os.Open(cfg.CustomMapTVDBPath)
		if err != nil {
			return err
		}

		err = readYamlFile(fTVDB, s)
		if err != nil {
			return err
		}

		cfg.TVDBMalMap = s
	}

	if cfg.CustomMapTMDB {
		am := &domain.AnimeMovies{}
		fTMDB, err := os.Open(cfg.CustomMapTMDBPath)
		if err != nil {
			return err
		}

		err = readYamlFile(fTMDB, am)
		if err != nil {
			return err
		}

		cfg.TMDBMalMap = am
	}

	return nil
}

func readYamlHTTP(resp *http.Response, mapping interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	err = yaml.Unmarshal(body, mapping)
	if err != nil {
		return err
	}

	return nil
}

func readYamlFile(f *os.File, mapping interface{}) error {
	defer f.Close()
	body, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, mapping)
	if err != nil {
		return err
	}

	return nil
}

func validateYaml(schema domain.CommunityMapUrls, yaml any) error {
	compiler := jsonschema.NewCompiler()
	sch, err := compiler.Compile(string(schema))
	if err != nil {
		return err
	}

	var v interface{}
	b, err := json.Marshal(yaml)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if err := sch.Validate(v); err != nil {
		return err
	}

	return nil
}
