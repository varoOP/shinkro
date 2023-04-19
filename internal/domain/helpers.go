package domain

import (
	"context"
)

func updateStart(ctx context.Context, s int) int {
	if s == 0 {
		return 1
	}
	return s
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
