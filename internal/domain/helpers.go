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
