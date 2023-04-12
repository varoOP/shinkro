package database

import (
	"github.com/varoOP/shinkuro/pkg/animelist"
	"github.com/varoOP/shinkuro/pkg/manami"
)

type Anime struct {
	AnidbID int
	Title   string
	MalID   int
	TvdbID  int
	TmdbID  int
}

type AnimeMap struct {
	Anime []Anime
}

func makeAnimeMap(m *manami.Manami, al *animelist.AnimeList) *AnimeMap {
	a := &AnimeMap{}
	anidbtotvmdb := al.AnidDbTvmDbmap()
	for _, v := range m.Data {
		malID, anidbID := v.GetID(`^https://myanimelist\.net/anime/([0-9]+)`), v.GetID(`^https://anidb\.net/anime/([0-9]+)`)
		tvdbID, tmdbID := al.GetTvmdbID(anidbID, anidbtotvmdb)
		if anidbID >= 0 {
			a.Anime = append(a.Anime, Anime{
				anidbID, v.Title, malID, tvdbID, tmdbID,
			})
		}
	}

	return a
}
