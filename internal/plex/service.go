package plex

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
)

type Service interface {
	Store(ctx context.Context, plex *domain.Plex) error
	// FindAll(ctx context.Context) ([]*domain.Plex, error)
	Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error)
	CheckPlex(plex *domain.Plex) bool
	// ProcessPlexScrobbleEvent(plex *domain.Plex) error
	// Delete(ctx context.Context, req *domain.DeletePlexRequest) error
}

type service struct {
	log            zerolog.Logger
	repo           domain.PlexRepo
	config         *domain.Config
	animeService   anime.Service
	mapService     mapping.Service
	malauthService malauth.Service
}

func NewService(log zerolog.Logger, config *domain.Config, repo domain.PlexRepo, animeSvc anime.Service, mapSvc mapping.Service, malauthSvc malauth.Service) Service {
	return &service{
		log:            log.With().Str("module", "plex").Logger(),
		repo:           repo,
		config:         config,
		animeService:   animeSvc,
		mapService:     mapSvc,
		malauthService: malauthSvc,
	}
}

func (s *service) Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error) {
	return s.repo.Get(ctx, req)
}

func (s *service) Store(ctx context.Context, plex *domain.Plex) error {
	return s.repo.Store(ctx, plex)
}

func (s *service) NewPlexClient(cfg *domain.Config) *domain.PlexClient {
	return &domain.PlexClient{
		Url:   cfg.PlexUrl,
		Token: cfg.PlexToken,
	}
}

func (s *service) CheckPlex(plex *domain.Plex) bool {
	if !isPlexUser(plex, s.config) {
		s.log.Debug().Err(errors.Wrap(errors.New("unauthorized plex user"), plex.Account.Title)).Msg("")
		return false
	}

	if !isEvent(plex) {
		s.log.Debug().Err(errors.Wrap(errors.New("plex event not supported"), string(plex.Event))).Msg("")
		return false
	}

	if !isAnimeLibrary(plex, s.config) {
		s.log.Debug().Err(errors.Wrap(errors.New("plex library not set as an anime library"), plex.Metadata.LibrarySectionTitle)).Msg("")
		return false
	}

	if !mediaType(plex) {
		s.log.Debug().Err(errors.Wrap(errors.New("plex media type not supported"), string(plex.Metadata.Type))).Msg("")
		return false
	}

	return true
}

func (s *service) ExtractSourceId(plex *domain.Plex) (string, int, error) {
	var (
		source string
		id     int
		err    error
	)
	// event := plex.GetPlexEvent()
	// if event == domain.PlexScrobbleEvent {
	// 	return nil
	// }

	agentAllowed, agent := isMetadataAgent(plex)
	if !agentAllowed {
		err = errors.Wrap(errors.New("metadata agent not supported"), string(agent))
		s.log.Debug().Err(err).Msg("")
		return "", 0, err
	}

	if agent == domain.HAMA || agent == domain.MALAgent {
		source, id, err = plex.Metadata.GUID.HamaMALAgent(agent)
		if err != nil {
			return "", 0, err
		}
	}

	if agent == domain.PlexAgent {
		if !isPlexClient(s.config) {
			err = errors.Wrap(errors.New("plex metadata agent cannot be used"), "Plex Token not set")
			s.log.Debug().Err(err).Msg("")
			return "", 0, err
		}

		guid := &domain.GUID{}
		if plex.Metadata.Type == domain.PlexEpisode {
			pc := s.NewPlexClient(s.config)
			g, err := pc.GetShowID(plex.Metadata.GrandparentKey)
			if err != nil {
				s.log.Debug().Err(err).Msg("")
				return "", 0, err
			}

			guid = g
		}

		source, id, err = guid.PlexAgent(plex.Metadata.Type)
		if err != nil {
			s.log.Debug().Err(err).Msg("")
			return "", 0, err
		}
	}

	return source, id, nil
}
