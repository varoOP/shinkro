package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/malauth"
	"github.com/varoOP/shinkuro/internal/mapping"
	"github.com/varoOP/shinkuro/internal/notifications"
	"github.com/varoOP/shinkuro/pkg/plex"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AnimeUpdate struct {
	client  *mal.Client
	db      *database.DB
	config  *config.Config
	anime   *mapping.Anime
	mapping *mapping.AnimeSeasonMap
	event   string
	inMap   bool
	media   *database.Media
	malid   int
	start   int
	rating  float32
	myList  *MyList
	malresp *mal.AnimeListStatus
}

type MyList struct {
	status     mal.AnimeStatus
	rewatchNum int
	epNum      int
	title      string
	picture    string
}

func NewAnimeUpdate(db *database.DB, cfg *config.Config) *AnimeUpdate {
	a := &AnimeUpdate{
		db:     db,
		config: cfg,
		malid:  -1,
		start:  -1,
	}

	return a
}

func (a *AnimeUpdate) SendUpdate(ctx context.Context) error {

	var err error
	c := malauth.NewOauth2Client(ctx, a.db)
	a.client = mal.NewClient(c)

	switch a.event {
	case "media.scrobble":
		if a.inMap {

			a.media.Ep = a.tvdbtoMal(ctx)
			err := a.updateWatchStatus(ctx)
			if err != nil {
				return err
			}
			return nil

		} else {
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
		}
	case "media.rate":

		if a.inMap {
			a.getStartID(ctx, a.anime.IsMultiSeason(ctx))
			err := a.updateRating(ctx)
			if err != nil {
				return err
			}
			return nil
		} else {
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
		}
	}
	return fmt.Errorf("%v - %v:not season 1 of anime, and not found in custom mapping", a.media.IdSource, a.media.Id)
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
		return nil
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

	aa, _, err := a.client.Anime.Details(ctx, a.malid, mal.Fields{"num_episodes", "title", "main_picture{medium}", "my_list_status{status,num_times_rewatched}"})
	if err != nil {
		return err
	}

	a.myList = &MyList{
		status:     aa.MyListStatus.Status,
		rewatchNum: aa.MyListStatus.NumTimesRewatched,
		epNum:      aa.NumEpisodes,
		title:      aa.Title,
		picture:    aa.MainPicture.Medium,
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

func (a *AnimeUpdate) createNotification() {
	d := notifications.NewDicord(a.config.Discord)
	if d.Url == "" {
		return
	}

	color := notifications.Color_watching
	if a.malresp.Status == mal.AnimeStatusCompleted {
		color = notifications.Color_completed
	}

	score := fmt.Sprintf("%v", a.malresp.Score)
	if score == "0" {
		score = "Not Scored"
	}

	d.Webhook = notifications.DiscordWebhook{
		Embeds: []notifications.Embeds{
			{
				Title: a.myList.title,
				URL:   fmt.Sprintf("https://myanimelist.net/anime/%v", a.malid),
				Color: color,
				Fields: []notifications.Fields{
					{
						Name:   "Status",
						Value:  cases.Title(language.Und).String(string(a.malresp.Status)),
						Inline: false,
					},
					{
						Name:   "Episodes Seen",
						Value:  fmt.Sprintf("%v / %v", a.malresp.NumEpisodesWatched, a.myList.epNum),
						Inline: true,
					},
					{
						Name:   "Your Score",
						Value:  score,
						Inline: true,
					},
				},
				Image: notifications.Image{
					URL: a.myList.picture,
				},
			},
		},
	}

	err := d.SendNotification()
	if err != nil {
		log.Println(err)
	}
}

func (a *AnimeUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pl := r.PostForm["payload"]
	ps := strings.Join(pl, "")

	if !isUserAgent(ps, a.config.User) {
		return
	}

	p, err := plex.NewPlexWebhook(ps)
	if err != nil {
		log.Println(err)
		return
	}

	if !isEvent(p.Event) {
		return
	}

	a.event = p.Event
	a.rating = p.Rating

	a.mapping, err = mapping.NewAnimeSeasonMap(a.config)
	if err != nil {
		log.Println("unable to load mapping", err)
		return
	}

	a.inMap, a.anime = a.mapping.CheckAnimeMap(p.Metadata.GrandparentTitle)

	a.media, err = database.NewMedia(p.Metadata.GUID, p.Metadata.Type)
	if err != nil {
		log.Println(err)
		return
	}

	err = a.SendUpdate(r.Context())
	if err != nil {
		log.Println(err)
		return
	}

	logUpdate(a.myList, a.malresp)
	a.createNotification()
	w.Write([]byte("Success"))
}
