package server

import (
	"context"
	"database/sql"
	"log"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/mapping"
	"github.com/varoOP/shinkuro/pkg/plex"
)

type AnimeUpdate struct {
	client *mal.Client
	db     *sql.DB
	anime  *mapping.Anime
	event  string
	inMap  bool
	show   *mapping.Show
	malid  int
	start  int
	rating float32
	myList *MyList
}

type AnimeCon struct {
	client  *mal.Client
	db      *sql.DB
	mapping *mapping.AnimeSeasonMap
}

type MyList struct {
	status     mal.AnimeStatus
	isRewatch  bool
	rewatchNum int
	epNum      int
	title      string
}

func NewAnimeCon(c *mal.Client, db *sql.DB) *AnimeCon {
	return &AnimeCon{
		client:  c,
		db:      db,
		mapping: nil,
	}
}

func NewAnimeUpdate(ctx context.Context, p *plex.PlexWebhook, ac *AnimeCon) (*AnimeUpdate, context.Context, error) {

	s, err := mapping.NewShow(ctx, p.Metadata.GUID)
	if err != nil {
		return nil, ctx, err
	}

	inMap, a := ac.mapping.CheckAnimeMap(ctx, p.Metadata.GrandparentTitle)

	am := &AnimeUpdate{
		client: ac.client,
		db:     ac.db,
		anime:  a,
		event:  p.Event,
		inMap:  inMap,
		show:   s,
		malid:  -1,
		start:  -1,
		rating: p.Rating,
	}

	return am, ctx, nil
}

func (am *AnimeUpdate) SendUpdate(ctx context.Context) error {

	var err error

	switch am.event {
	case "media.scrobble":
		if am.inMap {

			err = am.tvdbtoMal(ctx)
			if err != nil {
				return err
			}

		} else {
			if am.show.Season == 1 {
				am.malid, err = am.show.GetMalID(ctx, am.db)
				if err != nil {
					return err
				}
				err = am.updateWatchStatus(ctx)
				if err != nil {
					return err
				}
			}
		}
	case "media.rate":

		if am.inMap {
			am.getStartID(ctx, am.anime.IsMultiSeason(ctx))
			err = am.updateRating(ctx, am.rating)
			if err != nil {
				return err
			}
		} else {
			if am.show.Season == 1 {

				am.malid, err = am.show.GetMalID(ctx, am.db)
				if err != nil {
					return err
				}

				err = am.updateRating(ctx, am.rating)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (am *AnimeUpdate) tvdbtoMal(ctx context.Context) error {
	if !am.anime.IsMultiSeason(ctx) {
		am.getStartID(ctx, false)
		am.show.Ep = am.show.Ep - am.start + 1
		err := am.updateWatchStatus(ctx)
		if err != nil {
			return err
		}
	} else {
		am.getStartID(ctx, true)
		am.show.Ep = am.start + am.show.Ep - 1
		err := am.updateWatchStatus(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (am *AnimeUpdate) updateWatchStatus(ctx context.Context) error {

	options, err := am.newOptions(ctx)
	if err != nil {
		return err
	}

	l, _, err := am.client.Anime.UpdateMyListStatus(ctx, am.malid, options...)
	if err != nil {
		return err
	}

	logUpdate(am.myList, l)
	return nil
}

func (am *AnimeUpdate) newOptions(ctx context.Context) ([]mal.UpdateMyAnimeListStatusOption, error) {

	err := am.checkAnime(ctx)
	if err != nil {
		return nil, err
	}

	var options []mal.UpdateMyAnimeListStatusOption

	if am.myList.status == mal.AnimeStatusCompleted {
		if am.myList.epNum == am.show.Ep {
			am.myList.rewatchNum++
			options = append(options, mal.NumTimesRewatched(am.myList.rewatchNum))
			return options, nil
		} else {
			return nil, nil
		}
	}

	options = append(options, mal.NumEpisodesWatched(am.show.Ep))
	am.myList.status = mal.AnimeStatusWatching

	if am.myList.epNum == am.show.Ep {
		am.myList.status = mal.AnimeStatusCompleted
	}

	options = append(options, am.myList.status)
	return options, nil
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

func (am *AnimeUpdate) updateRating(ctx context.Context, r float32) error {

	err := am.checkAnime(ctx)
	if err != nil {
		return err
	}

	l, _, err := am.client.Anime.UpdateMyListStatus(ctx, am.malid, mal.Score(r))
	if err != nil {
		return err
	}

	logUpdate(am.myList, l)
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

func updateStart(ctx context.Context, s int) int {
	if s == 0 {
		return 1
	}
	return s
}

func logUpdate(ml *MyList, l *mal.AnimeListStatus) {

	log.Printf("%v - {Status:%v Score:%v Episodes_Watched:%v Rewatching:%v Times_Rewatched:%v}\n", ml.title, l.Status, l.Score, l.NumEpisodesWatched, l.IsRewatching, l.NumTimesRewatched)
}
