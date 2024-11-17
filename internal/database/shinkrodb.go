package database

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Anime struct {
	MainTitle    string `json:"title"`
	EnglishTitle string `json:"enTitle"`
	MalID        int    `json:"malid"`
	AnidbID      int    `json:"anidbid"`
	TvdbID       int    `json:"tvdbid"`
	TmdbID       int    `json:"tmdbid"`
	Type         string `json:"type"`
	ReleaseDate  string `json:"releaseDate"`
}

func getAnime() ([]Anime, error) {
	resp, err := http.Get("https://github.com/varoOP/shinkro-mapping/raw/main/shinkrodb/for-shinkro.json")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from shinkrodb")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	a := []Anime{}
	err = json.Unmarshal(body, &a)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal json")
	}

	return a, nil
}
