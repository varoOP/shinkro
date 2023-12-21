package domain

import (
	"context"

	"github.com/nstratos/go-myanimelist/mal"
)

type AnimeRepo interface {
	GetByID(ctx context.Context, req *GetAnimeRequest) (*Anime, error)
	StoreMultiple(anime []*Anime) error
}

type Anime struct {
	MALId           int                  `json:"malid"`
	MainTitle       string               `json:"title"`
	EnglishTitle    string               `json:"enTitle"`
	AniDBId         int                  `json:"anidbid"`
	TVDBId          int                  `json:"tvdbid"`
	TMDBId          int                  `json:"tmdbid"`
	AnimeType       string               `json:"type"`
	ReleaseDate     string               `json:"releaseDate"`
	EpisodeNum      int                  `json:"-"`
	AnimeListStatus *mal.AnimeListStatus `json:"-"`
	MyListDetails   ListDetails          `json:"-"`
}

type ListDetails struct {
	Status          mal.AnimeStatus
	RewatchNum      int
	TotalEpisodeNum int
	WatchedNum      int
	Title           string
	PictureURL      string
}

type IDTypes string

const TVDBID IDTypes = "tvdb_id"
const AniDBID IDTypes = "anidb_id"
const TMDBID IDTypes = "tmdb_id"

type AnimeUrl string

const ShinkroDB AnimeUrl = "https://github.com/varoOP/shinkrodb/raw/main/for-shinkro.json"

type GetAnimeRequest struct {
	IDtype IDTypes
	Id     int
}
