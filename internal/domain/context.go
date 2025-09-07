package domain

import (
	"context"
	"errors"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

// GetUserIDFromContext extracts userID from context
func GetUserIDFromContext(ctx context.Context) (int, error) {
	if userID, ok := ctx.Value(UserIDKey).(int); ok {
		return userID, nil
	}
	return 0, errors.New("userID not found in context")
}
