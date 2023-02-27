package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/animedb"
	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/pkg/plex"
)

func StartHttp(db *sql.DB, client *mal.Client, cfg *config.Config) {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test(w, r, db, client)
	})

	log.Println("Started listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}

func test(w http.ResponseWriter, r *http.Request, db *sql.DB, client *mal.Client) {

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

	if p.Event == "media.pause" || p.Event == "media.rate" {
		UpdateMal(r.Context(), &p, client, db)
	}

}

func UpdateMal(ctx context.Context, p *plex.PlexWebhook, client *mal.Client, db *sql.DB) {

	s := NewShow(p.Metadata.GUID)
	malid := s.GetMalID(db)
	title := p.Metadata.GrandparentTitle

	if s.Ep.Season == 1 && malid != 0 && p.Event == "media.scrobble" {

		status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.AnimeStatusWatching, mal.NumEpisodesWatched(s.Ep.No))
		log.Printf("%+v\n", *status)
		if err != nil {
			log.Println(err)
		}

	}

	if s.Ep.Season > 1 && malid != 0 && p.Event == "media.pause" {

		se := &animedb.SeasonMap{}
		se.GetSeasonMap()
		se.AddtoSeasonMap(ctx, malid, s.Ep.Season, title, client)
		/*
			a := fmt.Sprintf("Season_%v:", s.Ep.Season)

			for _, v := range se.Anime {
				if v.MalID == malid {
					for _, val := range v.Seasons {
						if strings.Contains(val, a) {
							f, _ := strings.CutPrefix(val, a)
							malid, err := strconv.Atoi(f)
							if err != nil {
								log.Fatalln(err)
							}

							status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.AnimeStatusWatching, mal.NumEpisodesWatched(s.Ep.No))
							log.Printf("%+v\n", *status)
							if err != nil {
								log.Println(err)
							}
							return
						}

					}
				}
			} */
	}

	if p.Event == "media.rate" && s.Ep.Season == 1 && malid != 0 {
		status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.Score(p.Rating))
		log.Printf("%+v\n", *status)
		if err != nil {
			log.Println(err)
		}
	}

	if s.Ep.Season > 1 && malid != 0 && p.Event == "media.rate" {
		se := &animedb.SeasonMap{}
		se.GetSeasonMap()

		a := fmt.Sprintf("Season_%v:", s.Ep.Season)

		for _, v := range se.Anime {
			if v.MalID == malid {
				for _, val := range v.Seasons {
					if strings.Contains(val, a) {
						f, _ := strings.CutPrefix(val, a)
						malid, err := strconv.Atoi(f)
						if err != nil {
							log.Fatalln(err)
						}

						status, _, err := client.Anime.UpdateMyListStatus(ctx, malid, mal.Score(p.Rating))
						log.Printf("%+v\n", *status)
						if err != nil {
							log.Println(err)
						}
						return
					}

				}
			}
		}
	}

}
