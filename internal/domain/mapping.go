package domain

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

const (
	communityMapTVDB = "https://github.com/varoOP/shinkro-mapping/raw/main/tvdb-mal.yaml"
	TVDBSchema       = "https://github.com/varoOP/shinkro-mapping/blob/main/.github/schema-tvdb.json"
	communityMapTMDB = "https://github.com/varoOP/shinkro-mapping/raw/main/tmdb-mal.yaml"
	TMDBSchema       = "https://github.com/varoOP/shinkro-mapping/raw/main/.github/schema-tmdb.json"
)

type AnimeSeasonMap struct {
	Anime []Anime `yaml:"anime" json:"anime"`
}

type Anime struct {
	Title    string    `yaml:"title" json:"title"`
	Synonyms []string  `yaml:"synonyms,omitempty" json:"synonyms,omitempty"`
	Seasons  []Seasons `yaml:"seasons" json:"seasons"`
}

type Seasons struct {
	Season int `yaml:"season" json:"season"`
	MalID  int `yaml:"mal-id" json:"mal-id"`
	Start  int `yaml:"start,omitempty" json:"start,omitempty"`
}

type AnimeMovies struct {
	AnimeMovie []AnimeMovie `yaml:"animeMovies" json:"animeMovies"`
}

type AnimeMovie struct {
	MainTitle string `yaml:"mainTitle" json:"mainTitle"`
	TMDBID    int    `yaml:"tmdbid" json:"tmdbid"`
	MALID     int    `yaml:"malid" json:"malid"`
}

func NewAnimeMaps(cfg *Config) (*AnimeSeasonMap, *AnimeMovies, error) {
	cfg.LocalMapsExist()
	err := loadCommunityMaps(cfg)
	if err != nil {
		return nil, nil, err
	}

	err = loadLocalMaps(cfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg.TVDBMalMap, cfg.TMDBMalMap, nil
}

func (a *Anime) IsMultiSeason(ctx context.Context, malid int) bool {
	var count int
	for _, s := range a.Seasons {
		if s.MalID == malid {
			count++
		}
	}

	return count > 1
}

func (s *AnimeSeasonMap) CheckMap(title string) (bool, *Anime) {
	for i, anime := range s.Anime {
		if title == anime.Title || synonymExists(anime.Synonyms, title) {
			return true, &Anime{
				Title:   s.Anime[i].Title,
				Seasons: s.Anime[i].Seasons,
			}
		}
	}

	return false, nil
}

func (am *AnimeMovies) CheckMap(tmdbid int) (bool, *AnimeMovie) {
	for _, animeMovie := range am.AnimeMovie {
		if animeMovie.TMDBID == tmdbid {
			return true, &animeMovie
		}
	}

	return false, nil
}

func loadCommunityMaps(cfg *Config) error {
	if !cfg.CustomMapTVDB {
		s := &AnimeSeasonMap{}
		respTVDB, err := http.Get(communityMapTVDB)
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
		am := &AnimeMovies{}
		respTMDB, err := http.Get(communityMapTMDB)
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

func loadLocalMaps(cfg *Config) error {
	if cfg.CustomMapTVDB {
		s := &AnimeSeasonMap{}
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
		am := &AnimeMovies{}
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

func synonymExists(s []string, title string) bool {

	for _, v := range s {
		if strings.EqualFold(v, title) {
			return true
		}
	}
	return false
}

func ChecklocalMaps(cfg *Config) (error, bool) {
	loadLocalMaps(cfg)
	localMapLoaded := false
	if cfg.CustomMapTVDB {
		if err := validateYaml(TVDBSchema, cfg.TVDBMalMap); err != nil {
			return err, false
		}

		localMapLoaded = true
	}

	if cfg.CustomMapTMDB {
		if err := validateYaml(TMDBSchema, cfg.TMDBMalMap); err != nil {
			return err, false
		}

		localMapLoaded = true
	}

	return nil, localMapLoaded
}

func validateYaml(schema string, yaml any) error {
	compiler := jsonschema.NewCompiler()
	sch, err := compiler.Compile(schema)
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
