package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/malauth"
)

type AnimeUpdate struct {
	Client  *mal.Client
	DB      *database.DB
	Config  *Config
	Anime   *Anime
	Mapping *AnimeSeasonMap
	Event   string
	InMap   bool
	Media   *database.Media
	Malid   int
	Start   int
	Ep      int
	Rating  float32
	MyList  *MyList
	Malresp *mal.AnimeListStatus
	Log     zerolog.Logger
	Notify  *Notification
}

type MyList struct {
	Status     mal.AnimeStatus
	RewatchNum int
	EpNum      int
	WatchedNum int
	Title      string
	Picture    string
}

func NewAnimeUpdate(db *database.DB, cfg *Config, log *zerolog.Logger, n *Notification) AnimeUpdate {
	return AnimeUpdate{
		DB:     db,
		Config: cfg,
		Malid:  -1,
		Start:  -1,
		Log:    log.With().Str("action", "animeUpdate").Logger(),
		Notify: n,
	}
}

func (a *AnimeUpdate) SendUpdate(ctx context.Context) error {
	c := malauth.NewOauth2Client(ctx, a.DB)
	a.Client = mal.NewClient(c)
	switch a.Event {
	case "media.scrobble":
		err := a.processScrobble(ctx)
		if err != nil {
			return err
		}

		return nil

	case "media.rate":
		err := a.processRate(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("unable to send update for %v (%v-%v)", a.Media.Title, a.Media.IdSource, a.Media.Id)
}

func (a *AnimeUpdate) processScrobble(ctx context.Context) error {
	var err error
	if a.InMap {
		err = a.tvdbtoMal(ctx)
		if err != nil {
			return err
		}

		err = a.updateWatchStatus(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	if a.Media.Season == 1 {
		a.Malid, err = a.Media.GetMalID(ctx, a.DB)
		if err != nil {
			return err
		}

		a.Ep = a.Media.Ep
		err = a.updateWatchStatus(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("unable to scrobble %v (%v-%v-%v), not found in database or custom map", a.Media.Title, a.Media.Type, a.Media.Agent, a.Media.Id)
}

func (a *AnimeUpdate) processRate(ctx context.Context) error {
	var err error
	if a.InMap {
		_, err := a.getStartID(ctx)
		if err != nil {
			return err
		}

		err = a.updateRating(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	if a.Media.Season == 1 {
		a.Malid, err = a.Media.GetMalID(ctx, a.DB)
		if err != nil {
			return err
		}

		err := a.updateRating(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("unable to rate %v (%v-%v-%v) not found in database or custom map", a.Media.Title, a.Media.Type, a.Media.Agent, a.Media.Id)
}

func (a *AnimeUpdate) tvdbtoMal(ctx context.Context) error {
	isMultiSeason, err := a.getStartID(ctx)
	if err != nil {
		return err
	}

	if isMultiSeason {
		a.Ep = a.Start + a.Media.Ep - 1
	} else {
		a.Ep = a.Media.Ep - a.Start + 1
	}

	if a.Ep <= 0 {
		return errors.New("episode calculated incorrectly")
	}

	return nil
}

func (a *AnimeUpdate) updateWatchStatus(ctx context.Context) error {
	options, complete, err := a.newOptions(ctx)
	if err != nil {
		return err
	}

	if complete {
		a.Log.Info().Msgf("%v is already marked complete on MAL", a.MyList.Title)
		return errors.New("complete")
	}

	l, _, err := a.Client.Anime.UpdateMyListStatus(ctx, a.Malid, options...)
	if err != nil {
		return err
	}

	a.Malresp = l
	return nil
}

func (a *AnimeUpdate) newOptions(ctx context.Context) ([]mal.UpdateMyAnimeListStatusOption, bool, error) {
	err := a.checkAnime(ctx)
	if err != nil {
		return nil, false, err
	}

	if a.Ep > a.MyList.EpNum && a.MyList.EpNum != 0 {
		return nil, true, fmt.Errorf("%v (%v-%v): anime in plex has more episodes for season than mal, modify custom mapping", a.Media.Title, a.Media.IdSource, a.Media.Id)
	}

	var options []mal.UpdateMyAnimeListStatusOption
	if a.MyList.Status == mal.AnimeStatusCompleted {
		if a.MyList.EpNum == a.Ep {
			a.MyList.RewatchNum++
			options = append(options, mal.NumTimesRewatched(a.MyList.RewatchNum))
			return options, false, nil
		}

		return nil, true, nil
	}

	if a.MyList.EpNum == a.Ep {
		a.MyList.Status = mal.AnimeStatusCompleted
		options = append(options, mal.FinishDate(time.Now().Local()))
	}

	if a.Ep == 1 && a.MyList.WatchedNum == 0 {
		options = append(options, mal.StartDate(time.Now().Local()))
	}

	if (a.Ep < a.MyList.EpNum || a.MyList.EpNum == 0) && a.Ep >= 1 {
		a.MyList.Status = mal.AnimeStatusWatching
	}

	options = append(options, mal.NumEpisodesWatched(a.Ep))
	options = append(options, a.MyList.Status)
	return options, false, nil
}

func (a *AnimeUpdate) checkAnime(ctx context.Context) error {
	aa, _, err := a.Client.Anime.Details(ctx, a.Malid, mal.Fields{"num_episodes", "title", "main_picture{medium,large}", "my_list_status{status,num_times_rewatched,num_episodes_watched}"})
	if err != nil {
		return err
	}

	picture := aa.MainPicture.Large
	if picture == "" {
		picture = aa.MainPicture.Medium
	}

	a.MyList = &MyList{
		Status:     aa.MyListStatus.Status,
		RewatchNum: aa.MyListStatus.NumTimesRewatched,
		EpNum:      aa.NumEpisodes,
		WatchedNum: aa.MyListStatus.NumEpisodesWatched,
		Title:      aa.Title,
		Picture:    picture,
	}

	return nil
}

func (a *AnimeUpdate) updateRating(ctx context.Context) error {
	err := a.checkAnime(ctx)
	if err != nil {
		return err
	}

	l, _, err := a.Client.Anime.UpdateMyListStatus(ctx, a.Malid, mal.Score(a.Rating))
	if err != nil {
		return err
	}

	a.Malresp = l
	return nil
}

func (a *AnimeUpdate) getStartID(ctx context.Context) (bool, error) {
	var isMultiSeason bool
	for _, anime := range a.Anime.Seasons {
		if a.Media.Season == anime.Season {
			if isMultiSeason = a.Anime.IsMultiSeason(ctx, anime.MalID); isMultiSeason {
				a.Malid = anime.MalID
				a.Start = anime.Start
			} else {
				if a.Media.Ep >= anime.Start {
					a.Malid = anime.MalID
					a.Start = anime.Start
				}
			}
		}
	}

	a.Start = updateStart(ctx, a.Start)
	if a.Malid <= 0 {
		return isMultiSeason, errors.New("no malid found")
	}

	return isMultiSeason, nil
}
