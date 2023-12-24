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

func (repo *AnimeUpdateRepo) Store(ctx context.Context, r *domain.AnimeUpdate) error {
	
	return nil
}

func (repo *AnimeUpdateRepo) GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error) {
	return nil, nil
}
