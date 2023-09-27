package domain

import (
	"context"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/pkg/plex"
)

type AnimeUpdate struct {
	Client      *mal.Client
	DB          *database.DB
	Config      *Config
	Plex        *plex.PlexWebhook
	Anime       *Anime
	AnimeMovie  *AnimeMovie
	TVDBMapping *AnimeTVDBMap
	TMDBMapping *AnimeMovies
	InTVDBMap   bool
	InTMDBMap   bool
	Media       *database.Media
	Malid       int
	Start       int
	Ep          int
	MyList      *MyList
	Malresp     *mal.AnimeListStatus
	Log         zerolog.Logger
	Notify      *Notification
}

type MyList struct {
	Status     mal.AnimeStatus
	RewatchNum int
	EpNum      int
	WatchedNum int
	Title      string
	Picture    string
}

type Key string

const (
	PlexPayload Key = "plexPayload"
	Agent       Key = "agent"
)

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
	if err := a.parseMedia(ctx); err != nil {
		return err
	}

	if err := a.getMapping(ctx); err != nil {
		return err
	}

	switch a.Plex.Event {
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

	return errors.Wrap(errors.Errorf("unable to send update for %v (%v-%v)", a.Media.Title, a.Media.IdSource, a.Media.Id), "plex event check failed")
}

func (a *AnimeUpdate) processScrobble(ctx context.Context) error {
	var err error
	if a.InTVDBMap {
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

	if a.InTMDBMap {
		a.Malid = a.AnimeMovie.MALID
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

	return errors.Wrap(errors.Errorf("unable to scrobble %v (%v-%v-%v)", a.Media.Title, a.Media.Type, a.Media.Agent, a.Media.Id), "not found in database or mapping")
}

func (a *AnimeUpdate) processRate(ctx context.Context) error {
	var err error
	if a.InTVDBMap {
		err := a.getStartID(ctx)
		if err != nil {
			return err
		}

		err = a.updateRating(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	if a.InTMDBMap {
		a.Malid = a.AnimeMovie.MALID
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

	return errors.Wrap(errors.Errorf("unable to rate %v (%v-%v-%v)", a.Media.Title, a.Media.Type, a.Media.Agent, a.Media.Id), "not found in database or mapping")
}

func (a *AnimeUpdate) tvdbtoMal(ctx context.Context) error {
	err := a.getStartID(ctx)
	if err != nil {
		return err
	}

	if a.Anime.UseMapping {
		a.Ep = a.Start + a.Media.Ep - 1
	} else {
		a.Ep = a.Media.Ep - a.Start + 1
	}

	if a.Ep <= 0 {
		return errors.Wrap(errors.New("episode calculated incorrectly"), "episode 0 or negative")
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
		return nil, true, errors.Wrap(errors.Errorf("%v (%v-%v): anime in plex has more episodes for season than mal", a.Media.Title, a.Media.IdSource, a.Media.Id), "update custom mappping to fix")
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

	l, _, err := a.Client.Anime.UpdateMyListStatus(ctx, a.Malid, mal.Score(a.Plex.Rating))
	if err != nil {
		return err
	}

	a.Malresp = l
	return nil
}

func (a *AnimeUpdate) getStartID(ctx context.Context) error {
	a.Malid = a.Anime.Malid
	a.Start = a.Anime.Start
	a.Start = updateStart(ctx, a.Start)
	return nil
}

func (a *AnimeUpdate) getMapping(ctx context.Context) error {
	var err error
	a.TVDBMapping, a.TMDBMapping, err = NewAnimeMaps(a.Config)
	if err != nil {
		return errors.Wrap(errors.New("unable to load custom mapping"), "check custom mapping against schema")
	}

	if a.Media.Type == "episode" {
		a.InTVDBMap, a.Anime = a.TVDBMapping.CheckMap(a.Media.Id, a.Media.Season, a.Media.Ep)
	}

	if a.Media.Type == "movie" {
		a.InTMDBMap, a.AnimeMovie = a.TMDBMapping.CheckMap(a.Media.Id)
	}

	return nil
}

func (a *AnimeUpdate) parseMedia(ctx context.Context) error {
	var (
		err     error
		pc      *plex.PlexClient
		usePlex bool = false
	)

	if a.Config.PlexToken != "" {
		pc = plex.NewPlexClient(a.Config.PlexUrl, a.Config.PlexToken)
		usePlex = true
	}

	a.Media, err = database.NewMedia(a.Plex, ctx.Value(Agent).(string), pc, usePlex)
	if err != nil {
		return err
	}

	return nil
}
