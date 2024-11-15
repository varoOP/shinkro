package domain

import (
	"context"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type MappingRepo interface {
	Store(ctx context.Context, m *MapSettings) error
	Get(ctx context.Context) (*MapSettings, error)
}

type Mapping struct {
	MapSettings MapSettings
	AnimeMap    AnimeMap
}

type MapSettings struct {
	TVDBEnabled       bool   `json:"tvdb_enabled"`
	TMDBEnabled       bool   `json:"tmdb_enabled"`
	CustomMapTVDBPath string `json:"tvdb_path"`
	CustomMapTMDBPath string `json:"tmdb_path"`
}

type AnimeMap struct {
	AnimeTVShows *AnimeTVShows
	AnimeMovies  *AnimeMovies
}

type AnimeMapDetails struct {
	Malid      int
	Start      int
	UseMapping bool
}

type AnimeTVShows struct {
	Anime []AnimeTV `yaml:"AnimeMap" json:"AnimeMap"`
}

type AnimeTV struct {
	Malid        int            `yaml:"malid" json:"malid"`
	Title        string         `yaml:"title" json:"title"`
	Type         string         `yaml:"type" json:"type"`
	Tvdbid       int            `yaml:"tvdbid" json:"tvdbid"`
	TvdbSeason   int            `yaml:"tvdbseason" json:"tvdbseason"`
	Start        int            `yaml:"start" json:"start"`
	UseMapping   bool           `yaml:"useMapping" json:"useMapping"`
	AnimeMapping []AnimeMapping `yaml:"animeMapping" json:"animeMapping"`
}

type AnimeMapping struct {
	TvdbSeason int `yaml:"tvdbseason" json:"tvdbseason"`
	Start      int `yaml:"start" json:"start"`
}

type AnimeMovies struct {
	AnimeMovie []AnimeMovie `yaml:"animeMovies" json:"animeMovies"`
}

type AnimeMovie struct {
	MainTitle string `yaml:"mainTitle" json:"mainTitle"`
	TMDBID    int    `yaml:"tmdbid" json:"tmdbid"`
	MALID     int    `yaml:"malid" json:"malid"`
}

type CommunityMapUrls string

const (
	CommunityMapTVDB CommunityMapUrls = "https://github.com/varoOP/shinkro-mapping/raw/main/tvdb-mal.yaml"
	CommunityMapTMDB CommunityMapUrls = "https://github.com/varoOP/shinkro-mapping/raw/main/tmdb-mal.yaml"
)

func (s *AnimeTVShows) CheckMap(tvdbid, tvdbseason, ep int) (bool, *AnimeTV) {
	candidates := s.findMatchingAnime(tvdbid, tvdbseason)
	if len(candidates) == 1 {
		return true, &candidates[0]
	} else if len(candidates) > 1 {
		anime := s.findBestMatchingAnime(ep, candidates)
		return true, &anime
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

func (s *AnimeTVShows) findMatchingAnime(tvdbid, tvdbseason int) []AnimeTV {
	var matchingAnime []AnimeTV
	for _, anime := range s.Anime {
		if tvdbid != anime.Tvdbid {
			continue
		}

		if !anime.UseMapping && tvdbseason == anime.TvdbSeason {
			matchingAnime = append(matchingAnime, anime)
			continue
		}

		matchingMappedAnime := s.findMatchingMappedAnime(anime, tvdbseason)
		if matchingMappedAnime != nil {
			return []AnimeTV{*matchingMappedAnime}
		}
	}

	return matchingAnime
}

func (s *AnimeTVShows) findMatchingMappedAnime(anime AnimeTV, tvdbseason int) *AnimeTV {
	if !anime.UseMapping {
		return nil
	}

	for _, animeMap := range anime.AnimeMapping {
		if tvdbseason == animeMap.TvdbSeason {
			anime.TvdbSeason = animeMap.TvdbSeason
			anime.Start = animeMap.Start
			return &anime
		}
	}

	return nil
}

func (s *AnimeTVShows) findBestMatchingAnime(ep int, candidates []AnimeTV) AnimeTV {
	var anime AnimeTV
	largestStart := 0
	for _, v := range candidates {
		if ep >= v.Start && v.Start >= largestStart {
			largestStart = v.Start
			anime = v
		}
	}

	return anime
}

func (ad *AnimeMapDetails) CalculateEpNum(oldEpNum int) int {
	if ad.UseMapping {
		return ad.Start + oldEpNum - 1
	}

	return oldEpNum - ad.Start + 1
}

func NewMapSettings(tvdb, tmdb bool, tvdbPath, tmdbPath string) *MapSettings {
	return &MapSettings{
		TVDBEnabled:       tvdb,
		TMDBEnabled:       tmdb,
		CustomMapTVDBPath: tvdbPath,
		CustomMapTMDBPath: tmdbPath,
	}
}

func (ms *MapSettings) LocalMapsExist() (bool, bool) {
	tvdb, tmdb := false, false

	if fileExists(ms.CustomMapTVDBPath) {
		tvdb = true
	}

	if fileExists(ms.CustomMapTMDBPath) {
		tmdb = true
	}

	return tvdb, tmdb
}

func (m *Mapping) LoadCommunityMaps(ctx context.Context, tvdb, tmdb bool) error {
	if !tvdb {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, string(CommunityMapTVDB), nil)
		if err != nil {
			return err
		}

		respTVDB, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		err = readYamlHTTP(respTVDB, m.AnimeMap.AnimeTVShows)
		if err != nil {
			return err
		}
	}

	if !tmdb {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, string(CommunityMapTMDB), nil)
		if err != nil {
			return err
		}

		respTMDB, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		err = readYamlHTTP(respTMDB, m.AnimeMap.AnimeMovies)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Mapping) LoadLocalMaps(tvdb, tmdb bool) error {
	if tvdb {
		fTVDB, err := os.Open(m.MapSettings.CustomMapTVDBPath)
		if err != nil {
			return err
		}

		err = readYamlFile(fTVDB, m.AnimeMap.AnimeTVShows)
		if err != nil {
			return err
		}
	}

	if tmdb {
		fTMDB, err := os.Open(m.MapSettings.CustomMapTMDBPath)
		if err != nil {
			return err
		}

		err = readYamlFile(fTMDB, m.AnimeMap.AnimeMovies)
		if err != nil {
			return err
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Open(path)
	return err == nil
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
