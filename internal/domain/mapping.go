package domain

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

const (
	communityMapTVDB = "https://github.com/varoOP/shinkro-mapping/raw/main/tvdb-mal.yaml"
	TVDBSchema       = "https://github.com/varoOP/shinkro-mapping/blob/main/.github/schema-tvdb.json"
	communityMapTMDB = "https://github.com/varoOP/shinkro-mapping/raw/main/tmdb-mal.yaml"
	TMDBSchema       = "https://github.com/varoOP/shinkro-mapping/raw/main/.github/schema-tmdb.json"
)

type AnimeTVDBMap struct {
	Anime []Anime `yaml:"AnimeMap"`
}

type Anime struct {
	Malid        int            `yaml:"malid"`
	Title        string         `yaml:"title"`
	Type         string         `yaml:"type"`
	Tvdbid       int            `yaml:"tvdbid"`
	TvdbSeason   int            `yaml:"tvdbseason"`
	Start        int            `yaml:"start"`
	UseMapping   bool           `yaml:"useMapping"`
	AnimeMapping []AnimeMapping `yaml:"animeMapping"`
}

type AnimeMapping struct {
	TvdbSeason int `yaml:"tvdbseason"`
	Start      int `yaml:"start"`
}

type AnimeMovies struct {
	AnimeMovie []AnimeMovie `yaml:"animeMovies" json:"animeMovies"`
}

type AnimeMovie struct {
	MainTitle string `yaml:"mainTitle" json:"mainTitle"`
	TMDBID    int    `yaml:"tmdbid" json:"tmdbid"`
	MALID     int    `yaml:"malid" json:"malid"`
}

func NewAnimeMaps(cfg *Config) (*AnimeTVDBMap, *AnimeMovies, error) {
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

func (s *AnimeTVDBMap) CheckMap(tvdbid, tvdbseason, ep int) (bool, *Anime) {
	candidates := s.findMatchingAnime(tvdbid, tvdbseason)
	if len(candidates) == 1 {
		return true, &candidates[0]
	} else if len(candidates) > 1 {
		anime := s.findBestMatchingAnime(ep, candidates)
		return true, &anime
	}

	return false, nil
}

func (s *AnimeTVDBMap) findMatchingAnime(tvdbid, tvdbseason int) []Anime {
	var matchingAnime []Anime
	for _, anime := range s.Anime {
		if tvdbid == anime.Tvdbid && tvdbseason == anime.TvdbSeason {
			matchingAnime = append(matchingAnime, anime)
		}

		if anime.UseMapping {
			for _, animeMap := range anime.AnimeMapping {
				if tvdbseason == animeMap.TvdbSeason {
					a := anime
					a.TvdbSeason = animeMap.TvdbSeason
					a.Start = animeMap.Start
					return []Anime{a}
				}
			}
		}
	}

	return matchingAnime
}

func (s *AnimeTVDBMap) findBestMatchingAnime(ep int, candidates []Anime) Anime {
	var anime Anime
	for _, v := range candidates {
		if ep >= v.Start {
			anime = v
		}
	}

	return anime
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
		s := &AnimeTVDBMap{}
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
		s := &AnimeTVDBMap{}
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
