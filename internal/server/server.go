package server

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/varoOP/shinkro/internal/mapping"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/notification"
	"github.com/varoOP/shinkro/internal/update"
)

type Server struct {
	log                 zerolog.Logger
	config              *domain.Config
	animeService        anime.Service
	mappingService      mapping.Service
	notificationService notification.Service
	lastUpdateNotified  string
}

func NewServer(log zerolog.Logger, config *domain.Config, animeSvc anime.Service, mappingSvc mapping.Service, notificationSvc notification.Service) *Server {
	return &Server{
		log:                 log.With().Str("module", "server").Logger(),
		config:              config,
		animeService:        animeSvc,
		mappingService:      mappingSvc,
		notificationService: notificationSvc,
	}
}

func (s *Server) Start() error {
	err := s.animeService.UpdateAnime()
	if err != nil {
		return err
	}

	if _, err := s.mappingService.Get(context.Background(), 1); errors.Is(err, sql.ErrNoRows) {
		_ = s.mappingService.Store(context.Background(), 1, &domain.MapSettings{
			TVDBEnabled:       false,
			TMDBEnabled:       false,
			CustomMapTMDBPath: "",
			CustomMapTVDBPath: "",
		})
	}

	c := cron.New(cron.WithLocation(time.UTC))
	_, err = c.AddFunc("0 1 * * MON", func() {
		s.animeService.UpdateAnime()
	})

	if err != nil {
		return err
	}

	c.Start()
	// Trigger update check after start and schedule a daily check
	go func() {
		time.Sleep(5 * time.Second)
		s.checkAndNotifyUpdate()
	}()
	_, _ = c.AddFunc("0 9 * * *", func() {
		s.checkAndNotifyUpdate()
	})
	return nil
}

// checkAndNotifyUpdate sends APP_UPDATE_AVAILABLE once per version in-memory
func (s *Server) checkAndNotifyUpdate() {
	if !s.config.CheckForUpdates {
		return
	}
	v := strings.ToLower(s.config.Version)
	if v == "" || strings.Contains(v, "dev") || strings.Contains(v, "nightly") {
		return
	}
	latest, err := update.LatestTag(context.Background())
	if err != nil || latest == "" {
		return
	}
	if !isUpdateAvailable(s.config.Version, latest) {
		return
	}
	if s.lastUpdateNotified == latest {
		return
	}
	s.notificationService.Send(domain.NotificationEventAppUpdateAvailable, domain.NotificationPayload{
		Subject:   "New update available!",
		Message:   latest,
		Event:     domain.NotificationEventAppUpdateAvailable,
		Timestamp: time.Now(),
	})
	s.lastUpdateNotified = latest
}

func isUpdateAvailable(current, latest string) bool {
	normalize := func(s string) []int {
		if len(s) > 0 && (s[0] == 'v' || s[0] == 'V') {
			s = s[1:]
		}
		parts := strings.Split(s, ".")
		out := make([]int, len(parts))
		for i, p := range parts {
			n := 0
			for _, ch := range p {
				if ch < '0' || ch > '9' {
					break
				}
				n = n*10 + int(ch-'0')
			}
			out[i] = n
		}
		return out
	}
	a := normalize(current)
	b := normalize(latest)
	for i := 0; i < len(a) || i < len(b); i++ {
		var x, y int
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		if x < y {
			return true
		}
		if x > y {
			return false
		}
	}
	return false
}
