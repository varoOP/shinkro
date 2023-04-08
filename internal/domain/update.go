package domain

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/malauth"
	"github.com/varoOP/shinkuro/internal/notification"
	"github.com/varoOP/shinkuro/pkg/plex"
)

type AnimeUpdate struct {
	client  *mal.Client
	db      *database.DB
	config  *Config
	anime   *Anime
	mapping *AnimeSeasonMap
	event   string
	inMap   bool
	media   *database.Media
	malid   int
	start   int
	rating  float32
	myList  *MyList
	malresp *mal.AnimeListStatus
	log     *zerolog.Logger
}

type MyList struct {
	status     mal.AnimeStatus
	rewatchNum int
	epNum      int
	title      string
	picture    string
}

func NewAnimeUpdate(db *database.DB, cfg *Config, log *zerolog.Logger) *AnimeUpdate {
	logger := log.With().Str("module", "domain").Logger()
	return &AnimeUpdate{
		db:     db,
		config: cfg,
		malid:  -1,
		start:  -1,
		log:    &logger,
	}
}

func (a *AnimeUpdate) SendUpdate(ctx context.Context) error {
	c := malauth.NewOauth2Client(ctx, a.db)
	a.client = mal.NewClient(c)

	switch a.event {
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

	return fmt.Errorf("%v - %v:not season 1 of anime, and not found in custom mapping", a.media.IdSource, a.media.Id)
}

func (a *AnimeUpdate) processScrobble(ctx context.Context) error {
	var err error
	if a.inMap {
		a.media.Ep = a.tvdbtoMal(ctx)
		err := a.updateWatchStatus(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	if a.media.Season == 1 {
		a.malid, err = a.media.GetMalID(ctx, a.db)
		if err != nil {
			return err
		}

		err := a.updateWatchStatus(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("unable to scrobble %v-%v-%v", a.media.Type, a.media.Agent, a.media.Id)
}

func (a *AnimeUpdate) processRate(ctx context.Context) error {
	var err error
	if a.inMap {
		a.getStartID(ctx, a.anime.IsMultiSeason(ctx))
		err := a.updateRating(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	if a.media.Season == 1 {
		a.malid, err = a.media.GetMalID(ctx, a.db)
		if err != nil {
			return err
		}

		err := a.updateRating(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("unable to rate %v-%v-%v", a.media.Type, a.media.Agent, a.media.Id)
}

func (a *AnimeUpdate) tvdbtoMal(ctx context.Context) int {
	if !a.anime.IsMultiSeason(ctx) {
		a.getStartID(ctx, false)
		ep := a.media.Ep - a.start + 1
		return ep
	} else {
		a.getStartID(ctx, true)
		ep := a.start + a.media.Ep - 1
		return ep
	}
}

func (a *AnimeUpdate) updateWatchStatus(ctx context.Context) error {

	options, complete, err := a.newOptions(ctx)
	if err != nil {
		return err
	}

	if complete {
		return fmt.Errorf("%v already marked complete on myanimelist", a.myList.title)
	}

	l, _, err := a.client.Anime.UpdateMyListStatus(ctx, a.malid, options...)
	if err != nil {
		return err
	}
	a.malresp = l
	return nil
}

func (a *AnimeUpdate) newOptions(ctx context.Context) ([]mal.UpdateMyAnimeListStatusOption, bool, error) {

	err := a.checkAnime(ctx)
	if err != nil {
		return nil, false, err
	}

	var options []mal.UpdateMyAnimeListStatusOption

	if a.myList.status == mal.AnimeStatusCompleted {
		if a.myList.epNum == a.media.Ep {
			a.myList.rewatchNum++
			options = append(options, mal.NumTimesRewatched(a.myList.rewatchNum))
			return options, false, nil
		} else if a.media.Ep > a.myList.epNum {
			return nil, true, fmt.Errorf("%v-%v: anime in plex has more episodes for season than mal, modify custom mapping", a.media.IdSource, a.media.Id)
		} else {
			return nil, true, nil
		}
	}

	options = append(options, mal.NumEpisodesWatched(a.media.Ep))
	a.myList.status = mal.AnimeStatusWatching

	if a.myList.epNum == a.media.Ep {
		a.myList.status = mal.AnimeStatusCompleted
		options = append(options, mal.FinishDate(time.Now().Local()))
	}

	if a.media.Ep == 1 {
		options = append(options, mal.StartDate(time.Now().Local()))
	}

	options = append(options, a.myList.status)
	return options, false, nil
}

func (a *AnimeUpdate) checkAnime(ctx context.Context) error {

	aa, _, err := a.client.Anime.Details(ctx, a.malid, mal.Fields{"num_episodes", "title", "main_picture{medium,large}", "my_list_status{status,num_times_rewatched}"})
	if err != nil {
		return err
	}

	picture := aa.MainPicture.Large
	if picture == "" {
		picture = aa.MainPicture.Medium
	}

	a.myList = &MyList{
		status:     aa.MyListStatus.Status,
		rewatchNum: aa.MyListStatus.NumTimesRewatched,
		epNum:      aa.NumEpisodes,
		title:      aa.Title,
		picture:    picture,
	}

	return nil
}

func (a *AnimeUpdate) updateRating(ctx context.Context) error {
	err := a.checkAnime(ctx)
	if err != nil {
		return err
	}

	l, _, err := a.client.Anime.UpdateMyListStatus(ctx, a.malid, mal.Score(a.rating))
	if err != nil {
		return err
	}

	a.malresp = l
	return nil
}

func (a *AnimeUpdate) getStartID(ctx context.Context, multi bool) {

	for _, anime := range a.anime.Seasons {
		if a.media.Season == anime.Season {
			if multi {
				a.malid = anime.MalID
				a.start = anime.Start
			} else {
				if a.media.Ep >= anime.Start {
					a.malid = anime.MalID
					a.start = anime.Start
				}
			}
		}
	}

	a.start = updateStart(ctx, a.start)
}

func (a *AnimeUpdate) createNotification(ctx context.Context) {
	d := notification.NewDicord(a.config.DiscordWebHookURL)
	if d.Url == "" {
		return
	}

	content := map[string]string{
		"event":           a.event,
		"title":           a.myList.title,
		"url":             fmt.Sprintf("https://myanimelist.net/anime/%v", a.malid),
		"status":          string(a.malresp.Status),
		"score":           strconv.Itoa(a.malresp.Score),
		"start_date":      string(a.malresp.StartDate),
		"finish_date":     string(a.malresp.FinishDate),
		"total_eps":       strconv.Itoa(a.myList.epNum),
		"watched_eps":     strconv.Itoa(a.malresp.NumEpisodesWatched),
		"times_rewatched": strconv.Itoa(a.malresp.NumTimesRewatched),
		"image_url":       a.myList.picture,
	}

	err := d.SendNotification(ctx, content)
	if err != nil {
		a.log.Debug().Err(err)
		return
	}
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
	a.log.Trace().Str("plexPayload", ps)

	if !isUserAgent(ps, a.config.PlexUser) {
		a.log.Debug().Msg("plex user or media's metadata agent did not match")
		return
	}

	p, err := plex.NewPlexWebhook(ps)
	if err != nil {
		a.log.Debug().Err(err).Msg("unable to unmarshal plex payload")
		return
	}

	if !isEvent(p.Event) {
		a.log.Trace().Str("event", p.Event).Msg("only accepting media.scrobble and media.rate events")
		return
	}

	a.event = p.Event
	a.rating = p.Rating

	a.mapping, err = NewAnimeSeasonMap(a.config)
	if err != nil {
		a.log.Info().Err(err).Msg("unable to load custom mapping")
		return
	}

	a.inMap, a.anime = a.mapping.CheckAnimeMap(p.Metadata.GrandparentTitle)

	a.media, err = database.NewMedia(p.Metadata.GUID, p.Metadata.Type)
	if err != nil {
		a.log.Info().Err(err).Msg("unable to parse media")
		return
	}

	err = a.SendUpdate(r.Context())
	if err != nil {
		a.log.Info().Err(err).Msg("failed to send update to myanimelist")
		return
	}

	a.log.Info().
		Str("status", string(a.malresp.Status)).
		Int("score", a.malresp.Score).
		Int("episdoesWatched", a.malresp.NumEpisodesWatched).
		Int("timesRewatched", a.malresp.NumTimesRewatched).
		Str("startDate", a.malresp.StartDate).
		Str("finishDate", a.malresp.FinishDate)

	a.createNotification(r.Context())
	w.Write([]byte("Success"))
}
