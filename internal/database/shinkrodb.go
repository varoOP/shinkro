package database

import (
	"encoding/json"
	"io"
	"net/http"
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
	resp, err := http.Get("https://github.com/varoOP/shinkrodb/raw/main/for-shinkro.json")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	a := []Anime{}
	err = json.Unmarshal(body, &a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
