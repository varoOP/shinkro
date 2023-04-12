package domain

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/malauth"
	"github.com/varoOP/shinkuro/pkg/plex"
)

type AnimeUpdate struct {
	client  *mal.Client
	db      *database.DB
	config  *Config
	anime   *Anime
	mapping *AnimeSeasonMap
	Event   string
	inMap   bool
	media   *database.Media
	Malid   int
	start   int
	ep      int
	rating  float32
	MyList  *MyList
	Malresp *mal.AnimeListStatus
	log     zerolog.Logger
	notify  *Notification
}

type MyList struct {
	Status     mal.AnimeStatus
	RewatchNum int
	EpNum      int
	WatchedNum int
	Title      string
	Picture    string
}

func NewAnimeUpdate(db *database.DB, cfg *Config, log *zerolog.Logger, n *Notification) *AnimeUpdate {
	return &AnimeUpdate{
		db:     db,
		config: cfg,
		Malid:  -1,
		start:  -1,
		log:    log.With().Str("action", "animeUpdate").Logger(),
		notify: n,
	}
}

func (a *AnimeUpdate) SendUpdate(ctx context.Context) error {
	c := malauth.NewOauth2Client(ctx, a.db)
	a.client = mal.NewClient(c)
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

	return fmt.Errorf("unable to send update for %v (%v-%v)", a.media.Title, a.media.IdSource, a.media.Id)
}

func (a *AnimeUpdate) processScrobble(ctx context.Context) error {
	var err error
	if a.inMap {
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

	if a.media.Season == 1 {
		a.Malid, err = a.media.GetMalID(ctx, a.db)
		if err != nil {
			return err
		}

		a.ep = a.media.Ep
		err = a.updateWatchStatus(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("unable to scrobble %v (%v-%v-%v), not found in database or custom map", a.media.Title, a.media.Type, a.media.Agent, a.media.Id)
}

func (a *AnimeUpdate) processRate(ctx context.Context) error {
	var err error
	if a.inMap {
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

	if a.media.Season == 1 {
		a.Malid, err = a.media.GetMalID(ctx, a.db)
		if err != nil {
			return err
		}

		err := a.updateRating(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("unable to rate %v (%v-%v-%v) not found in database or custom map", a.media.Title, a.media.Type, a.media.Agent, a.media.Id)
}

func (a *AnimeUpdate) tvdbtoMal(ctx context.Context) error {
	isMultiSeason, err := a.getStartID(ctx)
	if err != nil {
		return err
	}

	if isMultiSeason {
		a.ep = a.start + a.media.Ep - 1
	} else {
		a.ep = a.media.Ep - a.start + 1
	}

	if a.ep <= 0 {
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
		return fmt.Errorf("%v already marked complete on myanimelist", a.MyList.Title)
	}

	l, _, err := a.client.Anime.UpdateMyListStatus(ctx, a.Malid, options...)
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

	if a.ep > a.MyList.EpNum {
		return nil, true, fmt.Errorf("%v (%v-%v): anime in plex has more episodes for season than mal, modify custom mapping", a.media.Title, a.media.IdSource, a.media.Id)
	}

	var options []mal.UpdateMyAnimeListStatusOption
	if a.MyList.Status == mal.AnimeStatusCompleted {
		if a.MyList.EpNum == a.ep {
			a.MyList.RewatchNum++
			options = append(options, mal.NumTimesRewatched(a.MyList.RewatchNum))
			return options, false, nil
		}

		return nil, true, nil
	}

	if a.MyList.EpNum == a.ep {
		a.MyList.Status = mal.AnimeStatusCompleted
		options = append(options, mal.FinishDate(time.Now().Local()))
	}

	if a.ep == 1 && a.MyList.WatchedNum == 0 {
		options = append(options, mal.StartDate(time.Now().Local()))
	}

	if a.ep < a.MyList.EpNum && a.ep >= 1 {
		a.MyList.Status = mal.AnimeStatusWatching
	}

	options = append(options, mal.NumEpisodesWatched(a.ep))
	options = append(options, a.MyList.Status)
	return options, false, nil
}

func (a *AnimeUpdate) checkAnime(ctx context.Context) error {

	aa, _, err := a.client.Anime.Details(ctx, a.Malid, mal.Fields{"num_episodes", "title", "main_picture{medium,large}", "my_list_status{status,num_times_rewatched,num_episodes_watched}"})
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

	l, _, err := a.client.Anime.UpdateMyListStatus(ctx, a.Malid, mal.Score(a.rating))
	if err != nil {
		return err
	}

	a.Malresp = l
	return nil
}

func (a *AnimeUpdate) getStartID(ctx context.Context) (bool, error) {
	var isMultiSeason bool
	for _, anime := range a.anime.Seasons {
		if a.media.Season == anime.Season {
			if isMultiSeason = a.anime.IsMultiSeason(ctx, anime.MalID); isMultiSeason {
				a.Malid = anime.MalID
				a.start = anime.Start
			} else {
				if a.media.Ep >= anime.Start {
					a.Malid = anime.MalID
					a.start = anime.Start
				}
			}
		}
	}

	a.start = updateStart(ctx, a.start)
	if a.Malid <= 0 {
		return isMultiSeason, errors.New("no malid found")
	}

	return isMultiSeason, nil
}

func (a *AnimeUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.log.Trace().Msg("received bad request")
		return
	}

	pl := r.PostForm["payload"]
	ps := strings.Join(pl, "")
	a.log.Trace().Str("plexPayload", ps).Msg("received plex payload")

	if !isUserAgent(ps, a.config.PlexUser) {
		a.log.Debug().Msg("plex user or media's metadata agent did not match")
		return
	}

	p, err := plex.NewPlexWebhook(ps)
	if err != nil {
		a.log.Error().Err(err).Msg("unable to unmarshal plex payload")
		return
	}

	if !isEvent(p.Event) {
		a.log.Trace().Str("event", p.Event).Msg("only accepting media.scrobble and media.rate events")
		return
	}

	a.Event = p.Event
	a.rating = p.Rating

	a.mapping, err = NewAnimeSeasonMap(a.config)
	if err != nil {
		a.log.Error().Err(err).Msg("unable to load custom mapping")
		return
	}

	a.inMap, a.anime = a.mapping.CheckAnimeMap(p.Metadata.GrandparentTitle)

	a.media, err = database.NewMedia(p.Metadata.GUID, p.Metadata.Type, p.Metadata.GrandparentTitle)
	if err != nil {
		a.log.Error().Err(err).Msg("unable to parse media")
		return
	}

	err = a.SendUpdate(r.Context())
	notify(a, err)
	if err != nil {
		a.log.Error().Err(err).Msg("failed to send update to myanimelist")
		return
	}

	a.log.Info().
		Str("status", string(a.Malresp.Status)).
		Int("score", a.Malresp.Score).
		Int("episdoesWatched", a.Malresp.NumEpisodesWatched).
		Int("timesRewatched", a.Malresp.NumTimesRewatched).
		Str("startDate", a.Malresp.StartDate).
		Str("finishDate", a.Malresp.FinishDate).
		Msg("updated myanimelist successfully")

	w.Write([]byte("Success"))
}
