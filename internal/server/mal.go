package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/animedb"
	"github.com/varoOP/shinkuro/pkg/plex"
)

func UpdateMal(ctx context.Context, p *plex.PlexWebhook, c *mal.Client, db *sql.DB, se *animedb.SeasonMap) {

	var malid int
	s := NewShow(p.Metadata.GUID)
	title := p.Metadata.GrandparentTitle
	titleS := fmt.Sprintf("%v - %v", p.Metadata.GrandparentTitle, s.Ep.Season)
	inMap, a := getAnimeMap(title, se)

	switch p.Event {
	case "media.pause":
		if inMap {

			tvdbtoMal(ctx, a, c, s.Ep.Season, s.Ep.No, title)

		} else {
			if s.Ep.Season == 1 {

				malid = s.GetMalID(db)
				updateWatchStatus(ctx, c, title, malid, s.Ep.No)
			}
		}
	case "media.rate":
		if inMap {

			malid, _ := getStartID(a, s.Ep.Season, s.Ep.No, isMultiSeason(a))
			updateRating(ctx, c, titleS, malid, p.Rating)

		} else {
			if s.Ep.Season == 1 {

				malid = s.GetMalID(db)
				updateRating(ctx, c, titleS, malid, p.Rating)
			}
		}
	}
}

func getAnimeMap(title string, s *animedb.SeasonMap) (bool, *animedb.AnimeSeasons) {

	var inmap bool
	a := &animedb.AnimeSeasons{}

	for i, anime := range s.Anime {
		if title == anime.Title {

			inmap = true
			a.Title = s.Anime[i].Title
			a.Seasons = s.Anime[i].Seasons
			return inmap, a
		}
		inmap = false
	}
	return inmap, a
}

func tvdbtoMal(ctx context.Context, a *animedb.AnimeSeasons, c *mal.Client, season, ep int, title string) {

	if !isMultiSeason(a) {

		malid, start := getStartID(a, season, ep, false)

		updateWatchStatus(ctx, c, title, malid, ep-start+1)

	} else {

		malid, start := getStartID(a, season, ep, true)

		updateWatchStatus(ctx, c, title, malid, start+ep-1)

	}

}

func isMultiSeason(a *animedb.AnimeSeasons) bool {

	var count int
	var malid int

	for _, v := range a.Seasons {
		if v.Season == 1 {
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

func updateWatchStatus(ctx context.Context, c *mal.Client, title string, malid, ep int) {

	if !malID(malid, title) {
		return
	}

	status := mal.AnimeStatusWatching

	n, t := checkAnime(ctx, c, malid)

	title = t

	if n == ep {
		status = mal.AnimeStatusCompleted
	}

	l, _, err := c.Anime.UpdateMyListStatus(ctx, malid, status, mal.NumEpisodesWatched(ep))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%v - %+v\n", title, *l)
	}
}

func checkAnime(ctx context.Context, c *mal.Client, malid int) (int, string) {

	a, _, err := c.Anime.Details(ctx, malid, mal.Fields{"num_episodes", "title"})
	if err != nil {
		log.Println(err)
	} else {
		return a.NumEpisodes, a.Title
	}

	return 0, ""
}

func updateRating(ctx context.Context, c *mal.Client, title string, malid int, r float32) {

	if !malID(malid, title) {
		return
	}

	l, _, err := c.Anime.UpdateMyListStatus(ctx, malid, mal.Score(r))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%v - %+v\n", title, *l)
	}
}

func getStartID(a *animedb.AnimeSeasons, season, ep int, multi bool) (id, s int) {

	var malid int
	var start int

	for _, anime := range a.Seasons {
		if season == anime.Season {
			if multi {
				malid = anime.MalID
				start = anime.Start
			} else {
				if ep >= anime.Start {
					malid = anime.MalID
					start = anime.Start
				}
			}
		}
	}

	start = updateStart(start)

	return malid, start
}

func updateStart(s int) int {
	if s == 0 {
		return 1
	}
	return s
}

func malID(malid int, title string) bool {

	if malid == 0 {
		log.Printf("mal_id of %v not found. Update the mapping.\n", title)
		return false
	}

	return true
}
