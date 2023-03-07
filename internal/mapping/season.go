package mapping

import (
	"context"
	"io"
	"net/http"

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

func NewAnimeSeasonMap() (*AnimeSeasonMap, error) {
	s := &AnimeSeasonMap{}

	resp, err := http.Get("https://github.com/kyuketski/shinkuro-mapping/raw/main/tvdb-mal.yml")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(body, s)
	if err != nil {
		return nil, err
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

func (s *AnimeSeasonMap) CheckAnimeMap(ctx context.Context, title string) (bool, *Anime) {

	for i, anime := range s.Anime {
		if title == anime.Title || synonymExists(ctx, anime.Synonyms, title) {
			return true, &Anime{
				Title:   s.Anime[i].Title,
				Seasons: s.Anime[i].Seasons,
			}
		}
	}
	return false, nil
}

func synonymExists(ctx context.Context, s []string, title string) bool {

	for _, v := range s {
		if v == title {
			return true
		}
	}
	return false
}
