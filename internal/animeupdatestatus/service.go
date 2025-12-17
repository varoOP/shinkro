package animeupdatestatus

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	Store(ctx context.Context, status *domain.AnimeUpdateStatus) error
	GetByPlexID(ctx context.Context, plexID int64) (*domain.AnimeUpdateStatus, error)
	GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]domain.AnimeUpdateStatus, error)
}

type service struct {
	log  zerolog.Logger
	repo domain.AnimeUpdateStatusRepo
}

func NewService(log zerolog.Logger, repo domain.AnimeUpdateStatusRepo) Service {
	return &service{
		log:  log.With().Str("module", "animeupdatestatus").Logger(),
		repo: repo,
	}
}

func (s *service) Store(ctx context.Context, status *domain.AnimeUpdateStatus) error {
	return s.repo.Store(ctx, status)
}

func (s *service) GetByPlexID(ctx context.Context, plexID int64) (*domain.AnimeUpdateStatus, error) {
	return s.repo.GetByPlexID(ctx, plexID)
}

func (s *service) GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]domain.AnimeUpdateStatus, error) {
	return s.repo.GetByPlexIDs(ctx, plexIDs)
}

