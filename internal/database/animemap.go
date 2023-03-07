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
}

type AnimeMap struct {
	Anime []Anime
}

func makeAnimeMap(m *manami.Manami, al *animelist.AnimeList) *AnimeMap {

	a := &AnimeMap{}

	anidbtotvdb := al.AnidDbTvDbmap()

	for _, v := range m.Data {

		malID := v.SortToIDs(`^https://myanimelist\.net/anime/([0-9]+)`)
		anidbID := v.SortToIDs(`^https://anidb\.net/anime/([0-9]+)`)
		tvdbID := al.GetTvdbID(anidbID, anidbtotvdb)

		if anidbID != 0 {

			a.Anime = append(a.Anime, Anime{
				anidbID, v.Title, malID, tvdbID,
			})
		}
	}
	return a
}
