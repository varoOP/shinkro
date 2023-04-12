package domain

import (
	"context"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

const communityMapUrl = "https://github.com/varoOP/shinkuro-mapping/raw/main/tvdb-mal.yaml"

type AnimeSeasonMap struct {
	Anime []Anime `yaml:"anime"`
}

type Anime struct {
	Title    string    `yaml:"title"`
	Synonyms []string  `yaml:"synonyms,omitempty"`
	Seasons  []Seasons `yaml:"seasons"`
}

type Seasons struct {
	Season int `yaml:"season"`
	MalID  int `yaml:"mal-id"`
	Start  int `yaml:"start,omitempty"`
}

func NewAnimeSeasonMap(cfg *Config) (*AnimeSeasonMap, error) {
	s := &AnimeSeasonMap{}

	if cfg.CustomMapPath != "" {
		err := s.localMap(cfg.CustomMapPath)
		if err != nil {
			return nil, err
		}

		return s, nil
	}

	err := s.communityMap()
	if err != nil {
		return nil, err
	}

	return s, nil
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

func (s *AnimeSeasonMap) CheckAnimeMap(title string) (bool, *Anime) {

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

func (s *AnimeSeasonMap) communityMap() error {
	resp, err := http.Get(communityMapUrl)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = yaml.Unmarshal(body, s)
	if err != nil {
		return err
	}

	return nil
}

func (s *AnimeSeasonMap) localMap(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, s)
	if err != nil {
		return err
	}

	return nil
}

func synonymExists(s []string, title string) bool {

	for _, v := range s {
		if v == title {
			return true
		}
	}
	return false
}

func ChecklocalMap(path string) error {
	s := &AnimeSeasonMap{}
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()
	body, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, s)
	if err != nil {
		return err
	}

	return nil
}
