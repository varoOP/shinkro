package domain

import (
	"context"
	"time"
)

type AnimeUpdateStatusRepo interface {
	Store(ctx context.Context, status *AnimeUpdateStatus) error
	GetByPlexID(ctx context.Context, plexID int64) (*AnimeUpdateStatus, error)
	GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]AnimeUpdateStatus, error)
}

type AnimeUpdateStatus struct {
	ID          int64                    `json:"id"`
	PlexID      int64                    `json:"plexID"`
	MALID       int                      `json:"malID"`
	Status      AnimeUpdateStatusType    `json:"status"`
	ErrorType   AnimeUpdateErrorType     `json:"errorType,omitempty"`
	ErrorMessage string                  `json:"errorMessage,omitempty"`
	AnimeTitle   string                  `json:"animeTitle,omitempty"`
	SourceDB     PlexSupportedDBs        `json:"sourceDB,omitempty"`
	SourceID     int                      `json:"sourceID,omitempty"`
	SeasonNum    int                      `json:"seasonNum,omitempty"`
	EpisodeNum   int                      `json:"episodeNum,omitempty"`
	Timestamp    time.Time                `json:"timestamp"`
}

type AnimeUpdateStatusType string

const (
	AnimeUpdateStatusPending AnimeUpdateStatusType = "PENDING"
	AnimeUpdateStatusSuccess AnimeUpdateStatusType = "SUCCESS"
	AnimeUpdateStatusFailed  AnimeUpdateStatusType = "FAILED"
)

type AnimeUpdateErrorType string

const (
	AnimeUpdateErrorMALAuthFailed    AnimeUpdateErrorType = "MAL_AUTH_FAILED"
	AnimeUpdateErrorMappingNotFound  AnimeUpdateErrorType = "MAPPING_NOT_FOUND"
	AnimeUpdateErrorAnimeNotInDB     AnimeUpdateErrorType = "ANIME_NOT_IN_DB"
	AnimeUpdateErrorMALAPIFetchFailed AnimeUpdateErrorType = "MAL_API_FETCH_FAILED"
	AnimeUpdateErrorMALAPIUpdateFailed AnimeUpdateErrorType = "MAL_API_UPDATE_FAILED"
	AnimeUpdateErrorUnknown          AnimeUpdateErrorType = "UNKNOWN_ERROR"
)

