package mapping

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/varoOP/shinkuro/internal/config"
	"gopkg.in/yaml.v3"
)

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

func NewAnimeSeasonMap(cfg *config.Config) (*AnimeSeasonMap, error) {
	s := &AnimeSeasonMap{}

	if cfg.CustomMap {
		err := s.localMap(cfg.K.String("custom_map"))
		if err != nil {
			return nil, err
		}
	} else {
		err := s.communityMap()
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (a *Anime) IsMultiSeason(ctx context.Context) bool {

	var count, malid int

	for _, v := range a.Seasons {
		if v.Season == 1 || v.Season == 0 {
			malid = v.MalID
		}
		if malid == v.MalID {
			count++
		}

		if count > 1 {
			return true
		}
	}
	return false
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
	resp, err := http.Get("https://github.com/kyuketski/shinkuro-mapping/raw/main/tvdb-mal.yml")
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

func ChecklocalMap(path string) {

	s := &AnimeSeasonMap{}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal("error opening custom map", err)
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		log.Fatal("error reading custom map", err)
	}

	err = yaml.Unmarshal(body, s)
	if err != nil {
		log.Fatal(err)
	}

	s = nil
}
