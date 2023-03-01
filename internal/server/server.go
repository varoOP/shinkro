package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/animedb"
	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/pkg/plex"
)

func StartHttp(db *sql.DB, client *mal.Client, cfg *config.Config, se *animedb.SeasonMap) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		test(w, r, db, client, se)
	})

	log.Println("Started listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}

func test(w http.ResponseWriter, r *http.Request, db *sql.DB, client *mal.Client, se *animedb.SeasonMap) {

	p := &plex.PlexWebhook{}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal("Bad Request")
	}
	payload := r.PostForm["payload"]
	payloadstring := strings.Join(payload, "")

	if !strings.Contains(payloadstring, "com.plexapp.agents.hama") || !strings.Contains(payloadstring, "RudeusGreyrat") {
		return
	}

	err = json.Unmarshal([]byte(payloadstring), p)
	if err != nil {
		log.Println("Couldn't parse payload from Plex", err)
		return
	}

	if p.Event == "media.pause" || p.Event == "media.rate" {
		UpdateMal(r.Context(), p, client, db, se)
	}

}

func UpdateMal(ctx context.Context, p *plex.PlexWebhook, c *mal.Client, db *sql.DB, se *animedb.SeasonMap) {

	var malid int
	s := NewShow(p.Metadata.GUID)
	title := p.Metadata.GrandparentTitle
	inMap, a := getAnimeMap(title, se)

	switch p.Event {
	case "media.pause":
		if inMap {

			TvdbtoMal(ctx, a, c, s.Ep.Season, s.Ep.No, title)

		} else {
			if s.Ep.Season == 1 {

				malid = s.GetMalID(db)
				UpdateWatchStatus(ctx, c, title, malid, s.Ep.No)
			}
		}
	case "media.rate":
		if inMap {

			malid, _ := getStartID(a, s.Ep.Season, s.Ep.No, isMultiSeason(a))
			UpdateRating(ctx, c, title, malid, p.Rating)

		} else {
			if s.Ep.Season == 1 {

				malid = s.GetMalID(db)
				UpdateRating(ctx, c, title, malid, p.Rating)
			}
		}
	}
}

func getAnimeMap(title string, s *animedb.SeasonMap) (bool, *animedb.AnimeMT) {

	var inmap bool
	a := &animedb.AnimeMT{}

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

func TvdbtoMal(ctx context.Context, a *animedb.AnimeMT, c *mal.Client, season, ep int, title string) {

	if !isMultiSeason(a) {

		malid, start := getStartID(a, season, ep, false)

		UpdateWatchStatus(ctx, c, title, malid, ep-start+1)

	} else {

		malid, start := getStartID(a, season, ep, true)

		UpdateWatchStatus(ctx, c, title, malid, start+ep-1)

	}

}

func isMultiSeason(a *animedb.AnimeMT) bool {

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

/* func checkComplete (ctx context.Context, client *mal.Client, malid int) {

} */

func UpdateWatchStatus(ctx context.Context, c *mal.Client, title string, malid, ep int) {

	if !MalID(malid, title) {
		return
	}

	status, _, err := c.Anime.UpdateMyListStatus(ctx, malid, mal.AnimeStatusWatching, mal.NumEpisodesWatched(ep))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%v - %+v\n", title, *status)
	}
}

func UpdateRating(ctx context.Context, c *mal.Client, title string, malid int, r float32) {

	if !MalID(malid, title) {
		return
	}

	status, _, err := c.Anime.UpdateMyListStatus(ctx, malid, mal.Score(r))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%v - %+v\n", title, *status)
	}
}

func getStartID(a *animedb.AnimeMT, season, ep int, multi bool) (id, s int) {

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

func MalID(malid int, title string) bool {

	if malid == 0 {
		log.Printf("mal_id of %v not found. Update the mapping.\n", title)
		return false
	}

	return true
}
