package plexstatus

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	Store(ctx context.Context, status domain.PlexStatus) error
	GetByPlexID(ctx context.Context, plexID int64) (*domain.PlexStatus, error)
	GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]domain.PlexStatus, error)
	StoreSuccess(ctx context.Context, plex *domain.Plex) error
	StoreError(ctx context.Context, plex *domain.Plex, errorMsg string) error
}

type service struct {
	log  zerolog.Logger
	repo domain.PlexStatusRepo
}

func NewService(log zerolog.Logger, repo domain.PlexStatusRepo) Service {
	return &service{
		log:  log.With().Str("module", "plexStatus").Logger(),
		repo: repo,
	}
}

func (s *service) Store(ctx context.Context, status domain.PlexStatus) error {
	return s.repo.Store(ctx, status)
}

func (s *service) GetByPlexID(ctx context.Context, plexID int64) (*domain.PlexStatus, error) {
	return s.repo.GetByPlexID(ctx, plexID)
}

func (s *service) GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]domain.PlexStatus, error) {
	return s.repo.GetByPlexIDs(ctx, plexIDs)
}

func (s *service) StoreSuccess(ctx context.Context, plex *domain.Plex) error {
	title := plex.Metadata.GrandparentTitle

	if plex.Metadata.Type == "movie" {
		title = plex.Metadata.Title
	}

	status := domain.PlexStatus{
		Title:   title,
		Event:   string(plex.Event),
		Success: true,
		PlexID:  plex.ID,
	}
	return s.Store(ctx, status)
}

func (s *service) StoreError(ctx context.Context, plex *domain.Plex, errorMsg string) error {
	title := plex.Metadata.GrandparentTitle

	if plex.Metadata.Type == "movie" {
		title = plex.Metadata.Title
	}

	status := domain.PlexStatus{
		Title:    title,
		Event:    string(plex.Event),
		Success:  false,
		ErrorMsg: errorMsg,
		PlexID:   plex.ID,
	}
	return s.Store(ctx, status)
}
