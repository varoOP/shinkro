package plex

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	Store(ctx context.Context, plex *domain.Plex) error
	// FindAll(ctx context.Context) ([]*domain.Plex, error)
	Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error)
	// Delete(ctx context.Context, req *domain.DeletePlexRequest) error
}

type service struct {
	log  zerolog.Logger
	repo domain.PlexRepo
	//other services to come
}

func NewService(log zerolog.Logger, repo domain.PlexRepo) Service {
	return &service{
		log:  log.With().Str("module", "plex").Logger(),
		repo: repo,
	}
}

func (s *service) Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error) {
	return s.repo.Get(ctx, req)
}

func (s *service) Store(ctx context.Context, release *domain.Plex) error {
	return s.repo.Store(ctx, release)
}
