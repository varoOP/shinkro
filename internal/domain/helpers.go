package domain

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
	log.Printf("%v - {Status:%v Score:%v Episodes_Watched:%v Rewatching:%v Times_Rewatched:%v Start_Date:%v Finish_Date:%v}\n", ml.title, l.Status, l.Score, l.NumEpisodesWatched, l.IsRewatching, l.NumTimesRewatched, l.StartDate, l.FinishDate)
}

func isUserAgent(ps, user string) bool {
	if (strings.Contains(ps, "com.plexapp.agents.hama") || strings.Contains(ps, "net.fribbtastic.coding.plex.myanimelist")) && strings.Contains(ps, user) {
		return true
	}
	return false
}

func isEvent(e string) bool {
	if e == "media.rate" || e == "media.scrobble" {
		return true
	}
	return false
}
