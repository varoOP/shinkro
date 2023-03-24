package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/malauth"
	"github.com/varoOP/shinkuro/internal/mapping"
	"github.com/varoOP/shinkuro/pkg/plex"
)

type AnimeUpdate struct {
	client  *mal.Client
	db      *database.DB
	config  *config.Config
	anime   *mapping.Anime
	mapping *mapping.AnimeSeasonMap
	event   string
	inMap   bool
	show    *database.Show
	malid   int
	start   int
	ep      int
	rating  float32
	myList  *MyList
	malresp *mal.AnimeListStatus
}

type MyList struct {
	status     mal.AnimeStatus
	isRewatch  bool
	rewatchNum int
	epNum      int
	title      string
}

func NewAnimeUpdate(db *database.DB, cfg *config.Config) *AnimeUpdate {
	am := &AnimeUpdate{
		db:     db,
		config: cfg,
		malid:  -1,
		start:  -1,
	}

	return am
}

func (a *AnimeUpdate) SendUpdate(ctx context.Context) error {

	var err error
	c := malauth.NewOauth2Client(ctx, a.db)
	a.client = mal.NewClient(c)

	switch a.event {
	case "media.scrobble":
		if a.inMap {

			a.ep = a.tvdbtoMal(ctx)
			err := a.updateWatchStatus(ctx)
			if err != nil {
				return err
			}
			return nil

		} else {
			if a.show.Season == 1 {
				a.malid, err = a.show.GetMalID(ctx, a.db)
				if err != nil {
					return err
				}

				a.ep = a.show.Ep
				err := a.updateWatchStatus(ctx)
				if err != nil {
					return err
				}
				return nil
			}

			if a.malid > 0 {
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
			if a.show.Season == 1 {
				a.malid, err = a.show.GetMalID(ctx, a.db)
				if err != nil {
					return err
				}

				err := a.updateRating(ctx)
				if err != nil {
					return err
				}
				return nil
			}
			if a.malid > 0 {
				err := a.updateRating(ctx)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return fmt.Errorf("%v - %v:not season 1 of anime, and not found in custom mapping", a.show.IdSource, a.show.Id)
}

func (a *AnimeUpdate) tvdbtoMal(ctx context.Context) int {
	if !a.anime.IsMultiSeason(ctx) {
		a.getStartID(ctx, false)
		ep := a.show.Ep - a.start + 1
		return ep
	} else {
		a.getStartID(ctx, true)
		ep := a.start + a.show.Ep - 1
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
		if a.myList.epNum == a.ep {
			a.myList.rewatchNum++
			options = append(options, mal.NumTimesRewatched(a.myList.rewatchNum))
			return options, false, nil
		} else if a.ep > a.myList.epNum {
			return nil, true, fmt.Errorf("%v-%v: anime in plex has more episodes for season than mal, modify custom mapping", a.show.IdSource, a.show.Id)
		} else {
			return nil, true, nil
		}
	}

	options = append(options, mal.NumEpisodesWatched(a.ep))
	a.myList.status = mal.AnimeStatusWatching

	if a.myList.epNum == a.ep {
		a.myList.status = mal.AnimeStatusCompleted
	}

	options = append(options, a.myList.status)
	return options, false, nil
}

func (a *AnimeUpdate) checkAnime(ctx context.Context) error {

	aa, _, err := a.client.Anime.Details(ctx, a.malid, mal.Fields{"num_episodes", "title", "my_list_status"})
	if err != nil {
		return err
	}

	a.myList = &MyList{
		status:     aa.MyListStatus.Status,
		isRewatch:  aa.MyListStatus.IsRewatching,
		rewatchNum: aa.MyListStatus.NumTimesRewatched,
		epNum:      aa.NumEpisodes,
		title:      aa.Title,
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
		if a.show.Season == anime.Season {
			if multi {
				a.malid = anime.MalID
				a.start = anime.Start
			} else {
				if a.show.Ep >= anime.Start {
					a.malid = anime.MalID
					a.start = anime.Start
				}
			}
		}
	}

	a.start = updateStart(ctx, a.start)
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

	if p.Metadata.Type == "episode" {
		a.mapping, err = mapping.NewAnimeSeasonMap(a.config)
		if err != nil {
			log.Println("unable to load mapping", err)
			return
		}

		a.inMap, a.anime = a.mapping.CheckAnimeMap(p.Metadata.GrandparentTitle)

		a.show, err = database.NewShow(p.Metadata.GUID)
		if err != nil {
			log.Println(err)
			return
		}
	} else if p.Metadata.Type == "movie" {
		a.ep = 1
		a.malid, err = mapping.GetMovieMalID(p.Metadata.GUID)
		if err != nil {
			log.Println(err)
			return
		}
		a.show = &database.Show{
			IdSource: "",
			Id:       -1,
			Season:   -1,
			Ep:       -1,
		}
	} else {
		return
	}

	err = a.SendUpdate(r.Context())
	if err != nil {
		log.Println(err)
		return
	}

	logUpdate(a.myList, a.malresp)
	w.Write([]byte("Success"))
}
