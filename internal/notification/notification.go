package notification

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type AppNotification struct {
	Notification *domain.Notification
	Discord      *Discord
	log          zerolog.Logger
}

func NewAppNotification(url string, log *zerolog.Logger) *AppNotification {
	return &AppNotification{
		Notification: &domain.Notification{
			Error: make(chan error),
			Anime: make(chan domain.AnimeUpdate),
			Url:   url,
		},
		log: log.With().Str("module", "notification").Logger(),
	}
}

func NewNotificaitonPayload(a *domain.AnimeUpdate, err error) *domain.NotificationPayload {
	if a != nil {
		return &domain.NotificationPayload{
			// Event:          string(a.Plex.Event),
			// Title:          a.MyList.Title,
			// Url:            fmt.Sprintf("https://myanimelist.net/anime/%v", a.Malid),
			// Status:         string(a.Malresp.Status),
			// Score:          a.Malresp.Score,
			// StartDate:      a.Malresp.StartDate,
			// FinishDate:     a.Malresp.FinishDate,
			// TotalEps:       a.MyList.EpNum,
			// WatchedEps:     a.Malresp.NumEpisodesWatched,
			// TimesRewatched: a.Malresp.NumTimesRewatched,
			// ImageUrl:       a.MyList.Picture,
		}
	} else {
		return &domain.NotificationPayload{
			Title:   "Error",
			Message: fmt.Sprintf("`%v`", err),
		}
	}
}

func (n *AppNotification) ListenforNotification() {
	if n.Notification.Url == "" {
		return
	}

	for {
		select {
		case err := <-n.Notification.Error:
			n.log.Trace().Msg("received error in error channel")
			n.Notification.PayLoad = NewNotificaitonPayload(nil, err)
			n.CreateDiscord()
		case a := <-n.Notification.Anime:
			n.log.Trace().Msg("received notification in notification channel")
			n.Notification.PayLoad = NewNotificaitonPayload(&a, nil)
			n.CreateDiscord()
		}
	}
}

func (n *AppNotification) CreateDiscord() {
	n.Discord = NewDicord(n.Notification.Url, n.Notification.PayLoad, &n.log)
	n.log.Trace().Msg("built discord notification")
	if err := n.Discord.SendNotification(context.Background()); err != nil {
		n.log.Info().Err(err).Msg("unable to send discord notification")
	}
}
