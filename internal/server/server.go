package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/varoOP/shinkuro/pkg/plex"
)

func test(w http.ResponseWriter, r *http.Request, db *sql.DB) {

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

	if p.Event == "media.play" && strings.Contains(p.Metadata.GUID, "hama") {
		s := NewShow(p.Metadata.GUID)
		malid := s.GetMalID(db)
		fmt.Printf("%+v", s)
		fmt.Println("original id", s.Id, "db :", s.IdSource)
		fmt.Println("malid:", malid)
	}

	if p.Event == "media.rate" && p.Metadata.Type == "show" {
		s := NewShow(p.Metadata.GUID)
		malid := s.GetMalID(db)
		fmt.Printf("%+v", s)
		fmt.Println("malid:", malid)
	}

}

func StartHttp(db *sql.DB) {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test(w, r, db)
	})

	log.Fatal(http.ListenAndServe(":7011", nil))
}
