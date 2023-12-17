package database

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type AnimeUpdateRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewAnimeUpdateRepo(log zerolog.Logger, db *DB) domain.AnimeUpdateRepo {
	return &AnimeUpdateRepo{
		log: log.With().Str("repo", "animeupdate").Logger(),
		db:  db,
	}
}

func (r *AnimeUpdateRepo) Store(ctx context.Context, animeupdate *domain.AnimeUpdate) error {
	return nil
}

func (r *AnimeUpdateRepo) Delete(ctx context.Context, malid int) error {
	return nil
}

func (r *AnimeUpdateRepo) GetAnimeUpdates(ctx context.Context) ([]domain.AnimeUpdate, error) {
	return []domain.AnimeUpdate{}, nil
}
