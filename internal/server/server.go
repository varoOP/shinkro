package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/varoOP/shinkuro/pkg/plex"
)

func test(w http.ResponseWriter, r *http.Request, db *sql.DB, client *http.Client) {

	var p plex.PlexWebhook

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal("Bad Request")
	}
	payload := r.PostForm["payload"]
	payloadstring := strings.Join(payload, "")

	err = json.Unmarshal([]byte(payloadstring), &p)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if p.Event != "media.scrobble" || !strings.Contains(p.Metadata.GUID, "hama") || p.Account.Title != "RudeusGreyrat" {
		return
	}

	UpdateMal(&p, client, db)

}

func StartHttp(db *sql.DB, client *http.Client) {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test(w, r, db, client)
	})
	log.Println("Started listening on 7011.")
	log.Fatal(http.ListenAndServe(":7011", nil))
}
