package domain

import (
	"context"
	"strings"
)

func updateStart(ctx context.Context, s int) int {
	if s == 0 {
		return 1
	}
	return s
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

func notify(a *AnimeUpdate, err error) {
	if a.notify.Url == "" {
		return
	}

	if err != nil {
		a.notify.Error <- err
		return
	}

	a.notify.Anime <- *a
}
