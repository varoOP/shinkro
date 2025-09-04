package domain

import (
	"context"
	"time"
)

type NotificationRepo interface {
	List(ctx context.Context, userID int) ([]Notification, error)
	ListAll(ctx context.Context) ([]Notification, error)
	Find(ctx context.Context, params NotificationQueryParams) ([]Notification, int, error)
	FindByID(ctx context.Context, userID int, id int) (*Notification, error)
	Store(ctx context.Context, userID int, notification *Notification) error
	Update(ctx context.Context, userID int, notification *Notification) error
	Delete(ctx context.Context, userID int, notificationID int) error
}

type NotificationSender interface {
	Send(event NotificationEvent, payload NotificationPayload) error
	CanSend(event NotificationEvent) bool
	Name() string
}

type Notification struct {
	ID        int              `json:"id"`
	UserID    int              `json:"userID"`
	Name      string           `json:"name"`
	Type      NotificationType `json:"type"`
	Enabled   bool             `json:"enabled"`
	Events    []string         `json:"events"`
	Token     string           `json:"token"`
	APIKey    string           `json:"api_key"`
	Webhook   string           `json:"webhook"`
	Title     string           `json:"title"`
	Icon      string           `json:"icon"`
	Username  string           `json:"username"`
	Host      string           `json:"host"`
	Password  string           `json:"password"`
	Channel   string           `json:"channel"`
	Rooms     string           `json:"rooms"`
	Targets   string           `json:"targets"`
	Devices   string           `json:"devices"`
	Priority  int32            `json:"priority"`
	Topic     string           `json:"topic"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type NotificationPayload struct {
	Subject         string
	Message         string
	Event           NotificationEvent
	MediaName       string
	MALID           int
	AnimeLibrary    string
	EpisodesWatched int
	EpisodesTotal   int
	TimesRewatched  int
	PictureURL      string
	StartDate       string
	FinishDate      string
	AnimeStatus     string
	Score           int
	PlexEvent       PlexEvent
	PlexSource      PlexPayloadSource
	Timestamp       time.Time
	Sender          string
}

type NotificationType string

const (
	NotificationTypeDiscord NotificationType = "DISCORD"
	NotificationTypeGotify  NotificationType = "GOTIFY"
)

type NotificationEvent string

const (
	NotificationEventAppUpdateAvailable NotificationEvent = "APP_UPDATE_AVAILABLE"
	NotificationEventSuccess            NotificationEvent = "SUCCESS"
	NotificationEventError              NotificationEvent = "ERROR"
	NotificationEventTest               NotificationEvent = "TEST"
)

type NotificationEventArr []NotificationEvent

type NotificationQueryParams struct {
	UserID  int
	Limit   uint64
	Offset  uint64
	Sort    map[string]string
	Filters struct {
		Indexers   []string
		PushStatus string
	}
	Search string
}
