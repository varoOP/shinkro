package plex

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/varoOP/shinkro/internal/notification"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/animeupdate"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
	"github.com/varoOP/shinkro/internal/plexsettings"
	"github.com/varoOP/shinkro/internal/plexstatus"
)

type Service interface {
	Store(ctx context.Context, plex *domain.Plex) error
	Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error)
	ProcessPlex(ctx context.Context, plex *domain.Plex, agent *domain.PlexSupportedAgents) error
	GetPlexSettings(ctx context.Context) (*domain.PlexSettings, error)
	CountScrobbleEvents(ctx context.Context) (int, error)
	CountRateEvents(ctx context.Context) (int, error)
	GetRecent(ctx context.Context, limit int) ([]*domain.Plex, error)
}

type service struct {
	log                 zerolog.Logger
	repo                domain.PlexRepo
	plexettingsService  plexsettings.Service
	animeService        anime.Service
	mapService          mapping.Service
	malauthService      malauth.Service
	animeUpdateService  animeupdate.Service
	notificationService notification.Service
	plexStatusService   plexstatus.Service
}

func NewService(log zerolog.Logger, plexsettingsSvc plexsettings.Service, repo domain.PlexRepo, animeSvc anime.Service, mapSvc mapping.Service, malauthSvc malauth.Service, animeUpdateSvc animeupdate.Service, notificationSvc notification.Service, plexStatusSvc plexstatus.Service) Service {
	return &service{
		log:                 log.With().Str("module", "plex").Logger(),
		repo:                repo,
		plexettingsService:  plexsettingsSvc,
		animeService:        animeSvc,
		mapService:          mapSvc,
		malauthService:      malauthSvc,
		animeUpdateService:  animeUpdateSvc,
		notificationService: notificationSvc,
		plexStatusService:   plexStatusSvc,
	}
}

func (s *service) Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error) {
	return s.repo.Get(ctx, req)
}

func (s *service) Store(ctx context.Context, plex *domain.Plex) error {
	return s.repo.Store(ctx, plex)
}

func (s *service) GetPlexSettings(ctx context.Context) (*domain.PlexSettings, error) {
	return s.plexettingsService.Get(ctx)
}

func (s *service) ProcessPlex(ctx context.Context, plex *domain.Plex, agent *domain.PlexSupportedAgents) error {
	a, err := s.extractSourceIdForAnime(ctx, plex, agent)
	if err != nil {
		s.plexStatusService.StoreError(ctx, plex, err.Error())
		s.notificationService.Send(domain.NotificationEventError, domain.NotificationPayload{
			Message:      err.Error(),
			Subject:      "Failed to extract anime information",
			AnimeLibrary: plex.Metadata.LibrarySectionTitle,
			PlexEvent:    plex.Event,
			PlexSource:   plex.Source,
			Timestamp:    time.Now(),
		})
		return err
	}

	err = s.animeUpdateService.UpdateAnimeList(ctx, a, plex.Event)
	if err != nil {
		s.plexStatusService.StoreError(ctx, plex, err.Error())
		s.notificationService.Send(domain.NotificationEventError, domain.NotificationPayload{
			Message:      err.Error(),
			Subject:      "Failed to update MyAnimeList",
			AnimeLibrary: a.Plex.Metadata.LibrarySectionTitle,
			PlexEvent:    a.Plex.Event,
			PlexSource:   a.Plex.Source,
			Timestamp:    time.Now(),
		})

		return err
	}

	s.plexStatusService.StoreSuccess(ctx, plex)
	s.notificationService.Send(domain.NotificationEventSuccess, domain.NotificationPayload{
		MediaName:       a.ListDetails.Title,
		MALID:           a.MALId,
		AnimeLibrary:    a.Plex.Metadata.LibrarySectionTitle,
		EpisodesWatched: a.ListStatus.NumEpisodesWatched,
		EpisodesTotal:   a.ListDetails.TotalEpisodeNum,
		TimesRewatched:  a.ListStatus.NumTimesRewatched,
		PictureURL:      a.ListDetails.PictureURL,
		StartDate:       a.ListStatus.StartDate,
		FinishDate:      a.ListStatus.FinishDate,
		AnimeStatus:     string(a.ListStatus.Status),
		Score:           a.ListStatus.Score,
		PlexEvent:       a.Plex.Event,
		PlexSource:      a.Plex.Source,
		Timestamp:       time.Now(),
	})

	return nil
}

func (s *service) extractSourceIdForAnime(ctx context.Context, plex *domain.Plex, agent *domain.PlexSupportedAgents) (*domain.AnimeUpdate, error) {
	source, id, err := s.getSourceIDFromAgent(ctx, plex, agent)
	if err != nil {
		return nil, err
	}

	a := plex.SetAnimeFields(source, id)
	return &a, nil
}

func (s *service) getSourceIDFromAgent(ctx context.Context, p *domain.Plex, agent *domain.PlexSupportedAgents) (domain.PlexSupportedDBs, int, error) {
	switch *agent {
	case domain.HAMA, domain.MALAgent:
		return p.Metadata.GUID.HamaMALAgent(*agent)
	case domain.PlexAgent:
		return s.plexettingsService.HandlePlexAgent(ctx, p)
	}
	return "", 0, errors.New("unknown agent")
}

func (s *service) CountScrobbleEvents(ctx context.Context) (int, error) {
	return s.repo.CountScrobbleEvents(ctx)
}

func (s *service) CountRateEvents(ctx context.Context) (int, error) {
	return s.repo.CountRateEvents(ctx)
}

func (s *service) GetRecent(ctx context.Context, limit int) ([]*domain.Plex, error) {
	return s.repo.GetRecent(ctx, limit)
}
