package domain

type Notification struct {
	Error   chan error
	Anime   chan AnimeUpdate
	PayLoad *NotificationPayload
	Url     string
}

type NotificationPayload struct {
	Event          string
	Title          string
	Url            string
	Status         string
	Score          int
	StartDate      string
	FinishDate     string
	TotalEps       int
	WatchedEps     int
	TimesRewatched int
	ImageUrl       string
	Message        string
}
