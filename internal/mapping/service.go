package mapping

import (
	"context"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"gopkg.in/yaml.v3"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	NewMap(ctx context.Context, userID int) (*domain.AnimeMap, error)
	CheckForAnimeinMap(ctx context.Context, userID int, anime *domain.AnimeUpdate) (*domain.AnimeMapDetails, error)
	ValidateMap(ctx context.Context, yamlPath string, isTVDB bool) error
	Store(ctx context.Context, m *domain.MapSettings) error
	Get(ctx context.Context) (*domain.MapSettings, error)
}

type service struct {
	log  zerolog.Logger
	repo domain.MappingRepo

	mu        sync.RWMutex
	cachedMap *domain.AnimeMap
}

func NewService(log zerolog.Logger, repo domain.MappingRepo) Service {
	return &service{
		log:  log.With().Str("module", "mapping").Logger(),
		repo: repo,
	}
}

func (s *service) Store(ctx context.Context, m *domain.MapSettings) error {
	userID, err := domain.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	
	if err := s.repo.Store(ctx, userID, m); err != nil {
		return err
	}

	// Invalidate cache on store
	s.mu.Lock()
	s.cachedMap = nil
	s.mu.Unlock()

	return nil
}

func (s *service) Get(ctx context.Context) (*domain.MapSettings, error) {
	return s.repo.Get(ctx)
}

// NewMap bypasses cache and always loads fresh mappings
func (s *service) NewMap(ctx context.Context, userID int) (*domain.AnimeMap, error) {
	return s.loadMap(ctx, userID)
}

// loadMap encapsulates loading logic from repository
func (s *service) loadMap(ctx context.Context, userID int) (*domain.AnimeMap, error) {
	settings, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	mapping := &domain.Mapping{
		AnimeMap: &domain.AnimeMap{
			AnimeTVShows: &domain.AnimeTVShows{},
			AnimeMovies:  &domain.AnimeMovies{},
		},
		MapSettings: settings,
	}

	localTVDB, localTMDB := mapping.MapSettings.LocalMapsExist()
	if localTVDB || localTMDB {
		if err := mapping.LoadLocalMaps(localTVDB, localTMDB); err != nil {
			return nil, err
		}
	}
	if !localTVDB || !localTMDB {
		if err := mapping.LoadCommunityMaps(ctx, localTVDB, localTMDB); err != nil {
			return nil, err
		}
	}
	return mapping.AnimeMap, nil
}

// getCachedMap returns cached map or loads it if absent
func (s *service) getCachedMap(ctx context.Context, userID int) (*domain.AnimeMap, error) {
	s.mu.RLock()
	if s.cachedMap != nil {
		m := s.cachedMap
		s.mu.RUnlock()
		s.log.Debug().Msg("cache hit")
		return m, nil
	}
	s.mu.RUnlock()

	s.log.Debug().Msg("cache miss, loading map")
	m, err := s.loadMap(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cachedMap = m
	s.mu.Unlock()
	return m, nil
}

// reloadMap forces a refresh of the cache
func (s *service) reloadMap(ctx context.Context, userID int) (*domain.AnimeMap, error) {
	m, err := s.loadMap(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cachedMap = m
	s.mu.Unlock()
	return m, nil
}

func (s *service) CheckForAnimeinMap(ctx context.Context, userID int, anime *domain.AnimeUpdate) (*domain.AnimeMapDetails, error) {
	animeMap, err := s.getCachedMap(ctx, userID)
	if err != nil {
		return nil, err
	}

	if details, found := s.checkMap(animeMap, anime); found {
		return details, nil
	}

	animeMap, err = s.reloadMap(ctx, userID)
	if err != nil {
		return nil, err
	}

	s.log.Debug().Msg("reloaded map")

	if details, found := s.checkMap(animeMap, anime); found {
		return details, nil
	}

	return nil, errors.New("anime not found in map")
}

// checkMap encapsulates the lookup logic
func (s *service) checkMap(m *domain.AnimeMap, anime *domain.AnimeUpdate) (*domain.AnimeMapDetails, bool) {
	switch anime.Plex.Metadata.Type {
	case domain.PlexMovie:
		if inMap, movie := m.AnimeMovies.CheckMap(anime.SourceId); inMap {
			s.log.Debug().Msg("found anime movie in map")
			return &domain.AnimeMapDetails{Malid: movie.MALID, Start: 0}, true
		}

	case domain.PlexEpisode:
		if inMap, ep := m.AnimeTVShows.CheckMap(anime.SourceId, anime.SeasonNum, anime.EpisodeNum); inMap {
			s.log.Debug().Msg("found anime episode in map")
			return &domain.AnimeMapDetails{Malid: ep.Malid, Start: ep.Start, UseMapping: ep.UseMapping}, true
		}
	}
	return nil, false
}

func (s *service) ValidateMap(ctx context.Context, yamlPath string, isTVDB bool) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	var schemaURL domain.MapValidationUrls
	if isTVDB {
		schemaURL = domain.TVDBValidationUrl
	} else {
		schemaURL = domain.TMDBValidationUrl
	}

	raw, err := os.ReadFile(yamlPath)
	if err != nil {
		return errors.Errorf("could not read YAML file %q: %v", yamlPath, err)
	}

	var doc interface{}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return errors.Errorf("YAML parse error in %q: %v", yamlPath, err)
	}

	jsonData := toJSONCompatible(doc)
	compiler := jsonschema.NewCompiler()
	sch, err := compiler.Compile(string(schemaURL))
	if err != nil {
		return errors.Errorf("schema compile error for %s: %v", schemaURL, err)
	}

	if err := sch.Validate(jsonData); err != nil {
		s.log.Error().Err(err).Msgf("schema validation for %v custom map failed", yamlPath)
		return errors.Errorf("schema validation failed: %v", err)
	}

	return nil
}

func toJSONCompatible(v interface{}) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m2 := make(map[string]interface{}, len(x))
		for k, v2 := range x {
			m2[fmt.Sprint(k)] = toJSONCompatible(v2)
		}
		return m2
	case []interface{}:
		for i, u := range x {
			x[i] = toJSONCompatible(u)
		}
		return x
	default:
		return v
	}
}
