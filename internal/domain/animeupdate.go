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
}

type AnimeUpdate struct {
	ID          int64                `json:"id"`
	MALId       int                  `json:"malid"`
	SourceDB    PlexSupportedDBs     `json:"sourceDB"`
	SourceId    int                  `json:"sourceID"`
	EpisodeNum  int                  `json:"episodeNum"`
	SeasonNum   int                  `json:"seasonNum"`
	Timestamp   time.Time            `json:"timestamp"`
	ListDetails ListDetails          `json:"listDetails"`
	ListStatus  *mal.AnimeListStatus `json:"listStatus"`
	PlexId      int64                `json:"plexID"`
	Plex        *Plex                `json:"-"`
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

func (ap *AnimeUpdate) UpdateRating(ctx context.Context, client *mal.Client) error {
	err := ap.checkAnimeList(client, ctx)
	if err != nil {
		return err
	}

	l, _, err := client.Anime.UpdateMyListStatus(ctx, ap.MALId, mal.Score(ap.Plex.Rating))
	if err != nil {
		return err
	}

	ap.ListStatus = l
	return nil
}

func (ap *AnimeUpdate) UpdateWatchStatus(ctx context.Context, client *mal.Client) error {
	if err := ap.checkAnimeList(client, ctx); err != nil {
		return err
	}

	options, err := ap.newOptions(ctx)
	if err != nil {
		return err
	}

	l, _, err := client.Anime.UpdateMyListStatus(ctx, ap.MALId, options...)
	if err != nil {
		return err
	}

	ap.ListStatus = l
	return nil
}

func (ap *AnimeUpdate) newOptions(ctx context.Context) ([]mal.UpdateMyAnimeListStatusOption, error) {
	if err := ap.validateEpisodeNum(); err != nil {
		return nil, err
	}

	options, isComplete := ap.ListDetails.buildOptions(ap.EpisodeNum)
	if isComplete {
		return nil, errors.Errorf("%v is marked complete on MAL", ap.ListDetails.Title)
	}

	return options, nil
}

func (ap *AnimeUpdate) checkAnimeList(client *mal.Client, ctx context.Context) error {
	aa, _, err := client.Anime.Details(ctx, ap.MALId, mal.Fields{"num_episodes", "title", "main_picture{medium,large}", "my_list_status{status,num_times_rewatched,num_episodes_watched}"})
	if err != nil {
		return err
	}

	ap.ListDetails = ListDetails{
		Status:          aa.MyListStatus.Status,
		RewatchNum:      aa.MyListStatus.NumTimesRewatched,
		TotalEpisodeNum: aa.NumEpisodes,
		WatchedNum:      aa.MyListStatus.NumEpisodesWatched,
		Title:           aa.Title,
		PictureURL:      aa.MainPicture.Medium,
	}

	return nil
}

func (ap *AnimeUpdate) validateEpisodeNum() error {
	if ap.EpisodeNum > ap.ListDetails.TotalEpisodeNum && ap.ListDetails.TotalEpisodeNum != 0 {
		return errors.Errorf("number of episodes watched greater than total number of episodes: %v: Episode %v", ap.ListDetails.Title, ap.EpisodeNum)
	}
	return nil
}

func (ls *ListDetails) buildOptions(episodeNum int) ([]mal.UpdateMyAnimeListStatusOption, bool) {
	if ls.Status == mal.AnimeStatusCompleted && episodeNum != ls.TotalEpisodeNum {
		return nil, true
	}

	var options []mal.UpdateMyAnimeListStatusOption
	if ls.shouldIncrementRewatchNum(episodeNum) {
		ls.RewatchNum++
		options = append(options, mal.NumTimesRewatched(ls.RewatchNum))
	}

	if ls.isAnimeCompleted(episodeNum) {
		ls.Status = mal.AnimeStatusCompleted
		options = append(options, mal.FinishDate(time.Now().Local()))
	}

	if ls.isFirstEpisode(episodeNum) {
		options = append(options, mal.StartDate(time.Now().Local()))
	}

	if ls.isAnimeWatching(episodeNum) {
		ls.Status = mal.AnimeStatusWatching
	}

	options = append(options, mal.NumEpisodesWatched(episodeNum), ls.Status)
	return options, false
}

func (ls *ListDetails) shouldIncrementRewatchNum(episodeNum int) bool {
	return ls.Status == mal.AnimeStatusCompleted && ls.TotalEpisodeNum == episodeNum
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
