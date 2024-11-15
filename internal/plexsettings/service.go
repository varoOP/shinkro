package plexsettings

import (
	"context"
	"errors"
	"fmt"

	"github.com/dcarbone/zadapters/zstdlog"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/pkg/plex"
)

type Service interface {
	Store(ctx context.Context, ps domain.PlexSettings) (*domain.PlexSettings, error)
	Get(ctx context.Context) (*domain.PlexSettings, error)
	GetClient(ctx context.Context) (*plex.Client, error)
	HandlePlexAgent(ctx context.Context, p *domain.Plex) (domain.PlexSupportedDBs, int, error)
}

type service struct {
	log  zerolog.Logger
	repo domain.PlexSettingsRepo
}

func NewService(log zerolog.Logger, repo domain.PlexSettingsRepo) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) Store(ctx context.Context, ps domain.PlexSettings) (*domain.PlexSettings, error) {
	return s.repo.Store(ctx, ps)
}

func (s *service) Get(ctx context.Context) (*domain.PlexSettings, error) {
	return s.repo.Get(ctx)
}

func (s *service) GetClient(ctx context.Context) (*plex.Client, error) {
	ps, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !ps.PlexClientEnabled {
		return nil, errors.New("plex client disabled")
	}

	scheme := "http"
	if ps.TLS {
		scheme = "https"
	}

	c := plex.NewClient(plex.Config{
		Url:           fmt.Sprintf("%s://%s:%d", scheme, ps.Host, ps.Port),
		Token:         ps.Token,
		ClientID:      ps.ClientID,
		TLSSkipVerify: ps.TLSSkip,
		Log:           zstdlog.NewStdLoggerWithLevel(s.log.With().Str("client", "plex").Logger(), zerolog.TraceLevel),
	})

	return c, nil
}

func (s *service) HandlePlexAgent(ctx context.Context, p *domain.Plex) (domain.PlexSupportedDBs, int, error) {
	if p.Metadata.Type == domain.PlexEpisode {
		pc, err := s.GetClient(ctx)
		if err != nil {
			return "", 0, err
		}

		guid, err := pc.GetShowID(ctx, p.Metadata.GrandparentKey)
		if err != nil {
			return "", 0, err
		}

		id := domain.GUID{
			GUIDS: guid.GUIDS,
			GUID:  guid.GUID,
		}

		return id.PlexAgent(p.Metadata.Type)
	}
	return "", 0, nil
}
