package animeupdate

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	List(ctx context.Context) ([]domain.AnimeUpdate, error)
	Store(ctx context.Context, animeupdate *domain.AnimeUpdate) error
	Update(ctx context.Context, animeupdate *domain.AnimeUpdate) error
	Delete(ctx context.Context, malid int) error
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

func (s *service) List(ctx context.Context) ([]domain.AnimeUpdate, error) {
	return s.repo.GetAnimeUpdates(ctx)
}

func (s *service) Store(ctx context.Context, animeupdate *domain.AnimeUpdate) error {
	return s.repo.Store(ctx, animeupdate)
}

func (s *service) Update(ctx context.Context, animeupdate *domain.AnimeUpdate) error {
	return nil
}

func (s *service) Delete(ctx context.Context, malid int) error {
	return s.repo.Delete(ctx, malid)
}
