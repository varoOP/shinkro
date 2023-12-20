package plex

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	Store(ctx context.Context, plex *domain.Plex) error
	// FindAll(ctx context.Context) ([]*domain.Plex, error)
	Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error)
	CheckPlex(plex *domain.Plex) bool
	// Delete(ctx context.Context, req *domain.DeletePlexRequest) error
}

type service struct {
	log          zerolog.Logger
	repo         domain.PlexRepo
	config       *domain.Config
	animeService anime.Service
}

func NewService(log zerolog.Logger, config *domain.Config, repo domain.PlexRepo, animeSvc anime.Service) Service {
	return &service{
		log:          log.With().Str("module", "plex").Logger(),
		repo:         repo,
		config:       config,
		animeService: animeSvc,
	}
}

func (s *service) Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error) {
	return s.repo.Get(ctx, req)
}

func (s *service) Store(ctx context.Context, plex *domain.Plex) error {
	return s.repo.Store(ctx, plex)
}

func (s *service) CheckPlex(plex *domain.Plex) bool {
	if !isPlexUser(plex, s.config) {
		s.log.Debug().Err(errors.Wrap(errors.New("unauthorized plex user"), plex.Account.Title)).Msg("")
		return false
	}

	if !isEvent(plex) {
		s.log.Debug().Err(errors.Wrap(errors.New("plex event not supported"), plex.Event)).Msg("")
		return false
	}

	if !isAnimeLibrary(plex, s.config) {
		s.log.Debug().Err(errors.Wrap(errors.New("plex library not set as an anime library"), plex.Metadata.LibrarySectionTitle)).Msg("")
		return false
	}

	if !mediaType(plex) {
		s.log.Debug().Err(errors.Wrap(errors.New("plex media type not supported"), plex.Metadata.Type)).Msg("")
		return false
	}

	return true
}
