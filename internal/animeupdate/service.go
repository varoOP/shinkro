package animeupdate

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	Store(ctx context.Context, animeupdate *domain.AnimeUpdate) error
	GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error)
}

type service struct {
	log  zerolog.Logger
	repo domain.AnimeUpdateRepo
}

func NewService(log zerolog.Logger, repo domain.AnimeUpdateRepo) Service {
	return &service{
		log:  log.With().Str("module", "animeupdate").Logger(),
		repo: repo,
	}
}


func (s *service) Store(ctx context.Context, animeupdate *domain.AnimeUpdate) error {
	return s.repo.Store(ctx, animeupdate)
}

func (s *service) GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error) {
	return s.repo.GetByID(ctx, req)
}
