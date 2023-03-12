package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/mapping"
	"github.com/varoOP/shinkuro/pkg/plex"
)

func StartHttp(cfg *config.Config, ac *AnimeCon) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlePlexReq(w, r, ac, cfg)
	})

	log.Println("Started listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}

func handlePlexReq(w http.ResponseWriter, r *http.Request, ac *AnimeCon, cfg *config.Config) {

	p := &plex.PlexWebhook{}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pl := r.PostForm["payload"]
	ps := strings.Join(pl, "")

	if !isUserAgent(ps, cfg.User) {
		return
	}

	err = json.Unmarshal([]byte(ps), p)
	if err != nil {
		log.Println("Couldn't parse payload from Plex", err)
		return
	}

	ac.mapping, err = mapping.NewAnimeSeasonMap(cfg)
	if err != nil {
		log.Println("unable to load mapping", err)
		return
	}

	err = plexToMal(r.Context(), p, ac)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func plexToMal(ctx context.Context, p *plex.PlexWebhook, ac *AnimeCon) error {

	if !isEvent(p.Event) {
		return fmt.Errorf("incorrect event")
	}

	if p.Metadata.Type != "episode" {
		return fmt.Errorf("incorrect media type")
	}

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

func isUserAgent(ps, user string) bool {
	if !strings.Contains(ps, "com.plexapp.agents.hama") || !strings.Contains(ps, user) {
		return false
	}
	return true
}

func isEvent(e string) bool {
	events := []string{"media.rate", "media.scrobble"}
	for _, v := range events {
		if e == v {
			return true
		}
	}

	return false
}
