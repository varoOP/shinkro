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

func (ap *AnimeUpdate) UpdateWatchStatus(ctx context.Context, client *mal.Client) (bool, error) {
	if err := ap.checkAnimeList(client, ctx); err != nil {
		return false, err
	}

	options, isComplete, err := ap.newOptions(ctx)
	if err != nil {
		return false, err
	}

	if isComplete {
		return false, nil
	}

	l, _, err := client.Anime.UpdateMyListStatus(ctx, ap.MALId, options...)
	if err != nil {
		return false, err
	}

	ap.ListStatus = l
	return true, nil
}

func (ap *AnimeUpdate) newOptions(ctx context.Context) ([]mal.UpdateMyAnimeListStatusOption, bool, error) {
	if err := ap.validateEpisodeNum(); err != nil {
		return nil, false, err
	}

	options, isComplete := ap.ListDetails.buildOptions(ap.EpisodeNum)
	return options, isComplete, nil
}

func (ap *AnimeUpdate) checkAnimeList(client *mal.Client, ctx context.Context) error {
	aa, _, err := client.Anime.Details(ctx, ap.MALId, mal.Fields{"num_episodes", "title", "main_picture{medium,large}", "my_list_status{status,num_times_rewatched,num_episodes_watched}"})
	if err != nil {
		return err
	}

	picture := aa.MainPicture.Large
	if picture == "" {
		picture = aa.MainPicture.Medium
	}

	ap.ListDetails = ListDetails{
		Status:          aa.MyListStatus.Status,
		RewatchNum:      aa.MyListStatus.NumTimesRewatched,
		TotalEpisodeNum: aa.NumEpisodes,
		WatchedNum:      aa.MyListStatus.NumEpisodesWatched,
		Title:           aa.Title,
		PictureURL:      picture,
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
	var options []mal.UpdateMyAnimeListStatusOption
	var isComplete = false

	if ls.shouldIncrementRewatchNum(episodeNum) {
		ls.RewatchNum++
		options = append(options, mal.NumTimesRewatched(ls.RewatchNum))
	}

	if ls.isAnimeCompleted(episodeNum) {
		ls.Status = mal.AnimeStatusCompleted
		options = append(options, mal.FinishDate(time.Now().Local()))
		isComplete = true
	}

	if ls.isFirstEpisode(episodeNum) {
		options = append(options, mal.StartDate(time.Now().Local()))
	}

	if !isComplete && ls.isAnimeWatching(episodeNum) {
		ls.Status = mal.AnimeStatusWatching
	}

	options = append(options, mal.NumEpisodesWatched(episodeNum), ls.Status)
	return options, isComplete
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

// func (ap *AnimeUpdate) SetAnimeIDFromDB(source PlexSupportedDBs, id int) {
// 	switch source {
// 	case MAL:
// 		ap.MALId = id
// 	case TVDB:
// 		ap.TVDBId = id
// 	case TMDB:
// 		ap.TMDBId = id
// 	case AniDB:
// 		ap.AniDBId = id
// 	}
// }

// func NewAnimeUpdate(cfg *Config, log *zerolog.Logger, n *Notification) AnimeUpdate {
// 	return AnimeUpdate{
// 		Config: cfg,
// 		Malid:  -1,
// 		Start:  -1,
// 		Log:    log.With().Str("action", "animeUpdate").Logger(),
// 		Notify: n,
// 	}
// }

// func (a *AnimeUpdate) SendUpdate(ctx context.Context) error {
// 	c, err := malauth.NewOauth2Client(ctx, a.DB)
// 	if err != nil {
// 		return errors.Wrap(err, "unable to create new mal oauth2 client")
// 	}

// 	a.Client = mal.NewClient(c)
// 	if err := a.parseMedia(ctx); err != nil {
// 		return err
// 	}

// 	if err := a.getMapping(ctx); err != nil {
// 		return err
// 	}

// 	switch a.Plex.Event {
// 	case "media.scrobble":
// 		err := a.processScrobble(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		return nil

// 	case "media.rate":
// 		err := a.processRate(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}

// 	return errors.Wrap(errors.Errorf("unable to send update for %v (%v-%v)", a.Media.Title, a.Media.IdSource, a.Media.Id), "plex event check failed")
// }

// func (a *AnimeUpdate) processScrobble(ctx context.Context) error {
// 	var err error
// 	if a.InTVDBMap {
// 		err = a.tvdbtoMal(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		err = a.updateWatchStatus(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}

// 	if a.InTMDBMap {
// 		a.Malid = a.AnimeMovie.MALID
// 		err = a.updateWatchStatus(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}

// 	if a.Media.Season == 1 || a.Media.IdSource == "mal" {
// 		a.Malid, err = a.Media.GetMalID(ctx, a.DB)
// 		if err != nil {
// 			return err
// 		}

// 		a.Ep = a.Media.Ep
// 		err = a.updateWatchStatus(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}

// 	return errors.Wrap(errors.Errorf("unable to scrobble %v (%v-%v-%v)", a.Media.Title, a.Media.Type, a.Media.Agent, a.Media.Id), "not found in database or mapping")
// }

// func (a *AnimeUpdate) processRate(ctx context.Context) error {
// 	var err error
// 	if a.InTVDBMap {
// 		err := a.getStartID(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		err = a.updateRating(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}

// 	if a.InTMDBMap {
// 		a.Malid = a.AnimeMovie.MALID
// 		err = a.updateRating(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}

// 	if a.Media.Season == 1 || a.Media.IdSource == "mal" {
// 		a.Malid, err = a.Media.GetMalID(ctx, a.DB)
// 		if err != nil {
// 			return err
// 		}

// 		err := a.updateRating(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}

// 	return errors.Wrap(errors.Errorf("unable to rate %v (%v-%v-%v)", a.Media.Title, a.Media.Type, a.Media.Agent, a.Media.Id), "not found in database or mapping")
// }

// func (a *AnimeUpdate) tvdbtoMal(ctx context.Context) error {
// 	err := a.getStartID(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	if a.Anime.UseMapping {
// 		a.Ep = a.Start + a.Media.Ep - 1
// 	} else {
// 		a.Ep = a.Media.Ep - a.Start + 1
// 	}

// 	if a.Ep <= 0 {
// 		return errors.Wrap(errors.New("episode calculated incorrectly"), "episode 0 or negative")
// 	}

// 	return nil
// }

// func (a *AnimeUpdate) updateWatchStatus(ctx context.Context) error {
// 	options, complete, err := a.newOptions(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	if complete {
// 		a.Log.Info().Msgf("%v is already marked complete on MAL", a.MyList.Title)
// 		return errors.New("complete")
// 	}

// 	l, _, err := a.Client.Anime.UpdateMyListStatus(ctx, a.Malid, options...)
// 	if err != nil {
// 		return err
// 	}

// 	a.Malresp = l
// 	return nil
// }

// func (a *AnimeUpdate) newOptions(ctx context.Context) ([]mal.UpdateMyAnimeListStatusOption, bool, error) {
// 	err := a.checkAnime(ctx)
// 	if err != nil {
// 		return nil, false, err
// 	}

// 	if a.Ep > a.MyList.EpNum && a.MyList.EpNum != 0 {
// 		return nil, true, errors.Wrap(errors.Errorf("%v (%v-%v): anime in plex has more episodes for season than mal", a.Media.Title, a.Media.IdSource, a.Media.Id), "update custom mappping to fix")
// 	}

// 	var options []mal.UpdateMyAnimeListStatusOption
// 	if a.MyList.Status == mal.AnimeStatusCompleted {
// 		if a.MyList.EpNum == a.Ep {
// 			a.MyList.RewatchNum++
// 			options = append(options, mal.NumTimesRewatched(a.MyList.RewatchNum))
// 			return options, false, nil
// 		}

// 		return nil, true, nil
// 	}

// 	if a.MyList.EpNum == a.Ep {
// 		a.MyList.Status = mal.AnimeStatusCompleted
// 		options = append(options, mal.FinishDate(time.Now().Local()))
// 	}

// 	if a.Ep == 1 && a.MyList.WatchedNum == 0 {
// 		options = append(options, mal.StartDate(time.Now().Local()))
// 	}

// 	if (a.Ep < a.MyList.EpNum || a.MyList.EpNum == 0) && a.Ep >= 1 {
// 		a.MyList.Status = mal.AnimeStatusWatching
// 	}

// 	options = append(options, mal.NumEpisodesWatched(a.Ep))
// 	options = append(options, a.MyList.Status)
// 	return options, false, nil
// }

// func (a *AnimeUpdate) checkAnime(ctx context.Context) error {
// 	aa, _, err := a.Client.Anime.Details(ctx, a.Malid, mal.Fields{"num_episodes", "title", "main_picture{medium,large}", "my_list_status{status,num_times_rewatched,num_episodes_watched}"})
// 	if err != nil {
// 		return err
// 	}

// 	picture := aa.MainPicture.Large
// 	if picture == "" {
// 		picture = aa.MainPicture.Medium
// 	}

// 	a.MyList = &MyList{
// 		Status:     aa.MyListStatus.Status,
// 		RewatchNum: aa.MyListStatus.NumTimesRewatched,
// 		EpNum:      aa.NumEpisodes,
// 		WatchedNum: aa.MyListStatus.NumEpisodesWatched,
// 		Title:      aa.Title,
// 		Picture:    picture,
// 	}

// 	return nil
// }

// func (a *AnimeUpdate) updateRating(ctx context.Context) error {
// 	err := a.checkAnime(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	l, _, err := a.Client.Anime.UpdateMyListStatus(ctx, a.Malid, mal.Score(a.Plex.Rating))
// 	if err != nil {
// 		return err
// 	}

// 	a.Malresp = l
// 	return nil
// }

// func (a *AnimeUpdate) getStartID(ctx context.Context) error {
// 	a.Malid = a.Anime.Malid
// 	a.Start = a.Anime.Start
// 	// a.Start = updateStart(ctx, a.Start)
// 	return nil
// }

// func (a *AnimeUpdate) getMapping(ctx context.Context) error {
// 	var err error
// 	a.TVDBMapping, a.TMDBMapping, err = NewAnimeMaps(a.Config)
// 	if err != nil {
// 		return errors.Wrap(errors.New("unable to load custom mapping"), "check custom mapping against schema")
// 	}

// 	if a.Media.Type == "episode" {
// 		a.InTVDBMap, a.Anime = a.TVDBMapping.CheckMap(a.Media.Id, a.Media.Season, a.Media.Ep)
// 	}

// 	if a.Media.Type == "movie" {
// 		a.InTMDBMap, a.AnimeMovie = a.TMDBMapping.CheckMap(a.Media.Id)
// 	}

// 	return nil
// }

// func (a *AnimeUpdate) parseMedia(ctx context.Context) error {
// 	var (
// 		err     error
// 		pc      *plex.PlexClient
// 		usePlex bool = false
// 	)

// 	if a.Config.PlexToken != "" {
// 		pc = plex.NewPlexClient(a.Config.PlexUrl, a.Config.PlexToken)
// 		usePlex = true
// 	}

// 	a.Media, err = database.NewMedia(a.Plex, ctx.Value(Agent).(string), pc, usePlex)
// 	if err != nil {
// 		return err
// 	}

// 	a.Media.ConvertToTVDB(ctx, a.DB)
// 	return nil
// }
