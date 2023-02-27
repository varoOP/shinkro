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

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test(w, r, db, client, se)
	})

	log.Println("Started listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}

func test(w http.ResponseWriter, r *http.Request, db *sql.DB, client *mal.Client, se *animedb.SeasonMap) {

	var p plex.PlexWebhook

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal("Bad Request")
	}
	payload := r.PostForm["payload"]
	payloadstring := strings.Join(payload, "")

	if !strings.Contains(payloadstring, "com.plexapp.agents.hama") || !strings.Contains(payloadstring, "RudeusGreyrat") {
		return
	}

	err = json.Unmarshal([]byte(payloadstring), &p)
	if err != nil {
		log.Println("Couldn't parse payload from Plex", err)
		return
	}

	if p.Event == "media.scrobble" || p.Event == "media.rate" {
		UpdateMal(r.Context(), &p, client, db, se)
	}

}

func UpdateMal(ctx context.Context, p *plex.PlexWebhook, client *mal.Client, db *sql.DB, se *animedb.SeasonMap) {

	s := NewShow(p.Metadata.GUID)
	malid := s.GetMalID(db)
	title := p.Metadata.GrandparentTitle

	if s.Ep.Season == 1 && malid != 0 && p.Event == "media.scrobble" && !inSeasonMap(title, se) {

		status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.AnimeStatusWatching, mal.NumEpisodesWatched(s.Ep.No))
		log.Printf("%+v\n", *status)
		if err != nil {
			log.Println(err)
		}

	}

	if p.Event == "media.scrobble" && inSeasonMap(title, se) {

		a := getAnimeMT(title, se)

		tempUpdate(ctx, &a, client, s.Ep.Season, s.Ep.No)

	}

	if p.Event == "media.rate" && s.Ep.Season == 1 && malid != 0 && !inSeasonMap(title, se) {
		status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.Score(p.Rating))
		log.Printf("%+v\n", *status)
		if err != nil {
			log.Println(err)
		}
	}

	if p.Event == "media.rate" && inSeasonMap(title, se) {
		a := getAnimeMT(title, se)
		for _, v := range a.Seasons {
			if s.Ep.Season == v.Season {
				if s.Ep.No >= v.Start {
					malid = v.MalID
				}
			}
		}

		status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.Score(p.Rating))
		log.Printf("%+v\n", *status)
		if err != nil {
			log.Println(err)
		}
	}

}

func inSeasonMap(title string, s *animedb.SeasonMap) bool {
	for _, val := range s.Anime {
		if title == val.Title {
			return true
		}
	}
	return false
}

func getAnimeMT(title string, s *animedb.SeasonMap) animedb.AnimeMT {

	var a animedb.AnimeMT

	for i, val := range s.Anime {
		if title == val.Title {
			a = s.Anime[i]
		}
	}

	return a
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

func tempUpdate(ctx context.Context, a *animedb.AnimeMT, client *mal.Client, season, ep int) {

	var malid int
	var start int

	if !isMultiSeason(a) {

		for _, val := range a.Seasons {
			if season == val.Season {
				if ep >= val.Start {
					malid = val.MalID
					start = val.Start
				}
			}
		}

		if start == 0 {
			start++
		}

		status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.AnimeStatusWatching, mal.NumEpisodesWatched(ep-start+1))
		log.Printf("%+v\n", *status)
		if err != nil {
			log.Println(err)
		}

	} else {
		for _, val := range a.Seasons {
			if season == val.Season {
				malid = val.MalID
				start = val.Start
			}
		}

		if start == 0 {
			start++
		}

		status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.AnimeStatusWatching, mal.NumEpisodesWatched(start+ep-1))
		log.Printf("%+v\n", *status)
		if err != nil {
			log.Println(err)
		}

	}

}

/* func checkComplete (ctx context.Context, client *mal.Client, malid int) {

} */
