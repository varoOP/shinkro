package domain

import (
	"context"
	"time"
)

type PlexStatusRepo interface {
	Store(ctx context.Context, ps PlexStatus) error
	GetByPlexID(ctx context.Context, plexID int64) (*PlexStatus, error)
	GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]PlexStatus, error)
}

type PlexStatus struct {
	ID        int64            `json:"id"`
	Title     string           `json:"title"`
	Event     string           `json:"event"`
	Success   bool             `json:"success"`
	ErrorType PlexErrorType    `json:"errorType,omitempty"`
	ErrorMsg  string           `json:"errorMsg"`
	PlexID    int64            `json:"plexID"`
	TimeStamp time.Time        `json:"timestamp"`
}

type PlexErrorType string

const (
	PlexErrorAgentNotSupported PlexErrorType = "AGENT_NOT_SUPPORTED"
	PlexErrorExtractionFailed   PlexErrorType = "EXTRACTION_FAILED"
	PlexErrorUnknown            PlexErrorType = "UNKNOWN_ERROR"
)
