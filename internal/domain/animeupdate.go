package domain

import (
	"context"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/pkg/errors"
)

type AnimeUpdateRepo interface {
	Store(ctx context.Context, animeUpdate *AnimeUpdate) error
	GetByID(ctx context.Context, req *GetAnimeUpdateRequest) (*AnimeUpdate, error)
	Count(ctx context.Context) (int, error)
	GetRecentUnique(ctx context.Context, limit int) ([]*AnimeUpdate, error)
	GetByPlexID(ctx context.Context, plexID int64) (*AnimeUpdate, error)
	GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]*AnimeUpdate, error)
}

type AnimeUpdate struct {
	ID          int64               `json:"id"`
	MALId       int                 `json:"malid"`
	SourceDB    PlexSupportedDBs    `json:"sourceDB"`
	SourceId    int                 `json:"sourceID"`
	EpisodeNum  int                 `json:"episodeNum"`
	SeasonNum   int                 `json:"seasonNum"`
	Timestamp   time.Time           `json:"timestamp"`
	ListDetails ListDetails         `json:"listDetails"`
	ListStatus  mal.AnimeListStatus `json:"listStatus"`
	PlexId      int64               `json:"plexID"`
	Plex        *Plex               `json:"-"`
}

type ListDetails struct {
	Status          mal.AnimeStatus `json:"animeStatus"`
	RewatchNum      int             `json:"rewatchNum"`
	TotalEpisodeNum int             `json:"totalEpisodeNum"`
	WatchedNum      int             `json:"watchedNum"`
	Title           string          `json:"title"`
	PictureURL      string          `json:"pictureUrl"`
}

type Key string

const (
	PlexPayload Key = "plexPayload"
)

type GetAnimeUpdateRequest struct {
	Id int
}

// UpdateRatingWithStatus updates the rating using the provided list status from MAL API.
// The service layer should call MAL API first, then call this method with the result.
func (ap *AnimeUpdate) UpdateRatingWithStatus(status mal.AnimeListStatus) {
	ap.ListStatus = status
}

// BuildWatchStatusOptions builds the options for updating watch status.
// Returns the options that should be sent to MAL API.
func (ap *AnimeUpdate) BuildWatchStatusOptions() ([]mal.UpdateMyAnimeListStatusOption, error) {
	options, err := ap.newOptions()
	if err != nil {
		return nil, err
	}

	if len(options) == 0 {
		return nil, errors.New("no options to update")
	}

	return options, nil
}

// UpdateWatchStatusWithStatus updates the watch status using the provided list status from MAL API.
// The service layer should call MAL API first, then call this method with the result.
func (ap *AnimeUpdate) UpdateWatchStatusWithStatus(status mal.AnimeListStatus) {
	ap.ListStatus = status
}

// UpdateListDetails updates the list details from MAL API response data.
// This should be called by the service layer after fetching anime details from MAL.
func (ap *AnimeUpdate) UpdateListDetails(details ListDetails) {
	ap.ListDetails = details
}

func (ap *AnimeUpdate) newOptions() ([]mal.UpdateMyAnimeListStatusOption, error) {
	if err := ap.validateEpisodeNum(); err != nil {
		return nil, err
	}

	options := ap.ListDetails.buildOptions(ap.EpisodeNum)
	return options, nil
}

// BuildListDetailsFromMALResponse creates ListDetails from MAL API response.
// This is a pure transformation function - no I/O.
func BuildListDetailsFromMALResponse(status mal.AnimeStatus, rewatchNum, totalEpisodes, watchedNum int, title, pictureURL string) ListDetails {
	return ListDetails{
		Status:          status,
		RewatchNum:      rewatchNum,
		TotalEpisodeNum: totalEpisodes,
		WatchedNum:      watchedNum,
		Title:           title,
		PictureURL:      pictureURL,
	}
}

func (ap *AnimeUpdate) validateEpisodeNum() error {
	if ap.EpisodeNum > ap.ListDetails.TotalEpisodeNum && ap.ListDetails.TotalEpisodeNum != 0 {
		return errors.Errorf("number of episodes watched greater than total number of episodes: %v: Episode %v", ap.ListDetails.Title, ap.EpisodeNum)
	}
	return nil
}

func (ls *ListDetails) buildOptions(episodeNum int) []mal.UpdateMyAnimeListStatusOption {
	var options []mal.UpdateMyAnimeListStatusOption

	// Rewatching flow while status remains Completed
	if ls.Status == mal.AnimeStatusCompleted {
		if episodeNum < ls.TotalEpisodeNum || ls.TotalEpisodeNum == 0 {
			options = append(options, mal.IsRewatching(true), mal.NumEpisodesWatched(episodeNum), ls.Status)
			return options
		}

		if ls.isAnimeCompleted(episodeNum) {
				ls.RewatchNum++
				options = append(options, mal.NumTimesRewatched(ls.RewatchNum), mal.IsRewatching(false), mal.NumEpisodesWatched(episodeNum), ls.Status)
				return options
		}
	}

	// Normal progression before first completion
	if ls.isAnimeCompleted(episodeNum) {
		if ls.Status != mal.AnimeStatusCompleted {
			options = append(options, mal.FinishDate(time.Now().Local()))
		}
		ls.Status = mal.AnimeStatusCompleted
	}

	if ls.isFirstEpisode(episodeNum) {
		options = append(options, mal.StartDate(time.Now().Local()))
	}

	if ls.isAnimeWatching(episodeNum) {
		ls.Status = mal.AnimeStatusWatching
	}

	options = append(options, mal.NumEpisodesWatched(episodeNum), ls.Status)
	return options
}

func (ls *ListDetails) isAnimeCompleted(episodeNum int) bool {
	return ls.TotalEpisodeNum == episodeNum
}

func (ls *ListDetails) isFirstEpisode(episodeNum int) bool {
	return episodeNum == 1 && ls.WatchedNum == 0
}

func (ls *ListDetails) isAnimeWatching(episodeNum int) bool {
	return (episodeNum < ls.TotalEpisodeNum || ls.TotalEpisodeNum == 0) && episodeNum >= 1
}
