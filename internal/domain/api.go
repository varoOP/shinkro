package domain

import (
	"context"
	"time"
)

type APIRepo interface {
	Store(ctx context.Context, userID int, key *APIKey) error
	Delete(ctx context.Context, key string) error
	GetAllAPIKeys(ctx context.Context) ([]APIKey, error)
	GetKey(ctx context.Context, key string) (*APIKey, error)
	GetUserIDByAPIKey(ctx context.Context, key string) (int, error)
}

type APIKey struct {
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	Scopes    []string  `json:"scopes"`
	CreatedAt time.Time `json:"created_at"`
}
