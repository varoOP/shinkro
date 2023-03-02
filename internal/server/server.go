package server

import (
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
		log.Println("Bad", err)
		return
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

	if p.Event == "media.scrobble" || p.Event == "media.rate" {
		UpdateMal(r.Context(), p, client, db, se)
	}

}
