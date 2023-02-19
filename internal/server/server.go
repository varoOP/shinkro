package server

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/varoOP/shinkuro/pkg/plex"
	"golang.org/x/oauth2"
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

	if p.Event != "media.scrobble" || !strings.Contains(p.Metadata.GUID, "hama") {
		return
	}

	UpdateMal(&p, Oauth2Client, db)

}

func Authorize(w http.ResponseWriter, r *http.Request, cfg *oauth2.Config) {

	var (
		GrantType  oauth2.AuthCodeOption = oauth2.SetAuthURLParam("grant_type", "authorization_code")
		CodeVerify oauth2.AuthCodeOption = oauth2.SetAuthURLParam("code_verifier", Pkce)
	)

	r.ParseForm()

	s := r.Form.Get("state")
	if s != State {
		http.Error(w, "State invalid!", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found!", http.StatusBadRequest)
		return
	}

	token, err := cfg.Exchange(r.Context(), code, GrantType, CodeVerify)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	TokenCh <- token
	io.WriteString(w, "Token saved!\n")
}

func StartHttp(db *sql.DB, cfg *oauth2.Config) {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test(w, r, db)
	})

	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
		Authorize(w, r, cfg)
	})

	go func() {
		Oauth2Client = NewOauth2Client(cfg)
	}()

	log.Fatal(http.ListenAndServe(":7011", nil))
}
