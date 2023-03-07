package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/mapping"
	"github.com/varoOP/shinkuro/pkg/plex"
)

func StartHttp(cfg *config.Config, ac *AnimeCon) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		test(w, r, ac)
	})

	log.Println("Started listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}

func test(w http.ResponseWriter, r *http.Request, ac *AnimeCon) {

	p := &plex.PlexWebhook{}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("Bad", err)
		return
	}

	pl := r.PostForm["payload"]
	ps := strings.Join(pl, "")
	if !isUser(ps) {
		return
	}

	err = json.Unmarshal([]byte(ps), p)
	if err != nil {
		log.Println("Couldn't parse payload from Plex", err)
		return
	}

	err = plexToMal(r.Context(), p, ac)
	if err != nil {
		log.Println(err)
		return
	}
}

func plexToMal(ctx context.Context, p *plex.PlexWebhook, ac *AnimeCon) error {
	event := p.Event != "media.scrobble"
	if event || p.Event != "media.rate" {
		return nil
	}

	if p.Metadata.Type != "episode" {
		return nil
	}

	s, err := mapping.NewAnimeSeasonMap()
	if err != nil {
		return err
	}
	ac.mapping = s

	am, ctx, err := NewAnimeUpdate(ctx, p, ac)
	if err != nil {
		return err
	}

	err = am.SendUpdate(ctx)
	if err != nil {
		return err
	}

	return nil
}

func isUser(ps string) bool {
	if !strings.Contains(ps, "com.plexapp.agents.hama") || !strings.Contains(ps, "RudeusGreyrat") {
		return false
	}
	return true
}
