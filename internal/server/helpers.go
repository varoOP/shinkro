package server

import (
	"context"
	"log"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
)

func updateStart(ctx context.Context, s int) int {
	if s == 0 {
		return 1
	}
	return s
}

func logUpdate(ml *MyList, l *mal.AnimeListStatus) {

	log.Printf("%v - {Status:%v Score:%v Episodes_Watched:%v Rewatching:%v Times_Rewatched:%v}\n", ml.title, l.Status, l.Score, l.NumEpisodesWatched, l.IsRewatching, l.NumTimesRewatched)
}

func isUserAgent(ps, user string) bool {
	if !strings.Contains(ps, "com.plexapp.agents.hama") || !strings.Contains(ps, user) || !strings.Contains(ps, "net.fribbtastic.coding.plex.myanimelist") {
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
