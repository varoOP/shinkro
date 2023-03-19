package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/malauth"
	"github.com/varoOP/shinkuro/internal/mapping"
	"github.com/varoOP/shinkuro/pkg/plex"
)

type AnimeUpdate struct {
	client  *mal.Client
	db      *sql.DB
	config  *config.Config
	anime   *mapping.Anime
	mapping *mapping.AnimeSeasonMap
	event   string
	inMap   bool
	show    *mapping.Show
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

func NewAnimeUpdate(db *sql.DB, cfg *config.Config) *AnimeUpdate {
	am := &AnimeUpdate{
		db:     db,
		config: cfg,
		malid:  -1,
		start:  -1,
	}

	return am
}

func (am *AnimeUpdate) SendUpdate(ctx context.Context) error {

	var err error
	c := malauth.NewOauth2Client(ctx, am.db)
	am.client = mal.NewClient(c)

	switch am.event {
	case "media.scrobble":
		if am.inMap {

			am.ep = am.tvdbtoMal(ctx)
			err := am.updateWatchStatus(ctx)
			if err != nil {
				return err
			}
			return nil

		} else {
			if am.show.Season == 1 {
				am.malid, err = am.show.GetMalID(ctx, am.db)
				if err != nil {
					return err
				}
				err := am.updateWatchStatus(ctx)
				if err != nil {
					return err
				}
				return nil
			}

			if am.malid > 0 {
				err := am.updateWatchStatus(ctx)
				if err != nil {
					return err
				}
				return nil
			}
		}
	case "media.rate":

		if am.inMap {
			am.getStartID(ctx, am.anime.IsMultiSeason(ctx))
			err := am.updateRating(ctx)
			if err != nil {
				return err
			}
			return nil
		} else {
			if am.show.Season == 1 {
				am.malid, err = am.show.GetMalID(ctx, am.db)
				if err != nil {
					return err
				}

				err := am.updateRating(ctx)
				if err != nil {
					return err
				}
				return nil
			}
			if am.malid > 0 {
				err := am.updateRating(ctx)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return fmt.Errorf("%v - %v:not season 1 of anime, and not found in custom mapping", am.show.IdSource, am.show.Id)
}

func (am *AnimeUpdate) tvdbtoMal(ctx context.Context) int {
	if !am.anime.IsMultiSeason(ctx) {
		am.getStartID(ctx, false)
		ep := am.show.Ep - am.start + 1
		return ep
	} else {
		am.getStartID(ctx, true)
		ep := am.start + am.show.Ep - 1
		return ep
	}
}

func (am *AnimeUpdate) updateWatchStatus(ctx context.Context) error {

	options, complete, err := am.newOptions(ctx)
	if err != nil {
		return err
	}

	if complete {
		return nil
	}

	l, _, err := am.client.Anime.UpdateMyListStatus(ctx, am.malid, options...)
	if err != nil {
		return err
	}
	am.malresp = l
	return nil
}

func (am *AnimeUpdate) newOptions(ctx context.Context) ([]mal.UpdateMyAnimeListStatusOption, bool, error) {

	err := am.checkAnime(ctx)
	if err != nil {
		return nil, false, err
	}

	var options []mal.UpdateMyAnimeListStatusOption

	if am.myList.status == mal.AnimeStatusCompleted {
		if am.myList.epNum == am.ep {
			am.myList.rewatchNum++
			options = append(options, mal.NumTimesRewatched(am.myList.rewatchNum))
			return options, false, nil
		} else if am.ep > am.myList.epNum {
			return nil, true, fmt.Errorf("%v-%v: anime in plex has more episodes for season than mal, modify custom mapping", am.show.IdSource, am.show.Id)
		} else {
			return nil, true, nil
		}
	}

	options = append(options, mal.NumEpisodesWatched(am.ep))
	am.myList.status = mal.AnimeStatusWatching

	if am.myList.epNum == am.ep {
		am.myList.status = mal.AnimeStatusCompleted
	}

	options = append(options, am.myList.status)
	return options, false, nil
}

func (am *AnimeUpdate) checkAnime(ctx context.Context) error {

	a, _, err := am.client.Anime.Details(ctx, am.malid, mal.Fields{"num_episodes", "title", "my_list_status"})
	if err != nil {
		return err
	}

	am.myList = &MyList{
		status:     a.MyListStatus.Status,
		isRewatch:  a.MyListStatus.IsRewatching,
		rewatchNum: a.MyListStatus.NumTimesRewatched,
		epNum:      a.NumEpisodes,
		title:      a.Title,
	}

	return nil
}

func (am *AnimeUpdate) updateRating(ctx context.Context) error {

	err := am.checkAnime(ctx)
	if err != nil {
		return err
	}

	l, _, err := am.client.Anime.UpdateMyListStatus(ctx, am.malid, mal.Score(am.rating))
	if err != nil {
		return err
	}
	am.malresp = l
	return nil
}

func (am *AnimeUpdate) getStartID(ctx context.Context, multi bool) {

	for _, anime := range am.anime.Seasons {
		if am.show.Season == anime.Season {
			if multi {
				am.malid = anime.MalID
				am.start = anime.Start
			} else {
				if am.show.Ep >= anime.Start {
					am.malid = anime.MalID
					am.start = anime.Start
				}
			}
		}
	}

	am.start = updateStart(ctx, am.start)
}

func (a *AnimeUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := &plex.PlexWebhook{}

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

	err = json.Unmarshal([]byte(ps), p)
	if err != nil {
		log.Println("Couldn't parse payload from Plex", err)
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

		a.show, err = mapping.NewShow(p.Metadata.GUID)
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
		a.show = &mapping.Show{
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
