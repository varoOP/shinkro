package server

import (
	"context"
	"database/sql"
	"log"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/animedb"
	"github.com/varoOP/shinkuro/pkg/plex"
)

type AnimeUpdate struct {
	client       *mal.Client
	db           *sql.DB
	animeSeasons *animedb.AnimeSeasons
	title        string
	event        string
	inMap        bool
	show         *Show
	malid        int
	start        int
	rating       float32
}

type MyList struct {
	status     mal.AnimeStatus
	isRewatch  bool
	rewatchNum int
	epNum      int
	title      string
}

func NewAnimeUpdate(ctx context.Context, p *plex.PlexWebhook, c *mal.Client, db *sql.DB, sm *animedb.SeasonMap) (*AnimeUpdate, context.Context, error) {

	am := &AnimeUpdate{}

	s, err := NewShow(ctx, p.Metadata.GUID)
	if err != nil {
		return am, ctx, err
	}

	inMap, a := checkAnimeMap(ctx, p.Metadata.GrandparentTitle, sm)

	am = &AnimeUpdate{
		c,
		db,
		a,
		p.Metadata.GrandparentTitle,
		p.Event,
		inMap,
		s,
		-1,
		-1,
		p.Rating,
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
			if am.show.season == 1 {
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
			am.getStartID(ctx, am.animeSeasons.IsMultiSeason(ctx))
			err = am.updateRating(ctx, am.rating)
			if err != nil {
				return err
			}
		} else {
			if am.show.season == 1 {

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

func checkAnimeMap(ctx context.Context, title string, s *animedb.SeasonMap) (bool, *animedb.AnimeSeasons) {

	var inmap bool
	a := &animedb.AnimeSeasons{}

	for i, anime := range s.Anime {
		if title == anime.Title || synonymExists(ctx, anime.Synonyms, title) {

			inmap = true
			a.Title = s.Anime[i].Title
			a.Seasons = s.Anime[i].Seasons
			return inmap, a
		}
		inmap = false
	}
	return inmap, a
}

func (am *AnimeUpdate) tvdbtoMal(ctx context.Context) error {

	if !am.animeSeasons.IsMultiSeason(ctx) {

		am.getStartID(ctx, false)
		am.show.ep = am.show.ep - am.start + 1
		err := am.updateWatchStatus(ctx)
		if err != nil {
			return err
		}

	} else {

		am.getStartID(ctx, true)
		am.show.ep = am.start + am.show.ep - 1
		err := am.updateWatchStatus(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (am *AnimeUpdate) updateWatchStatus(ctx context.Context) error {

	ml, err := am.checkAnime(ctx)
	if err != nil {
		return err
	}

	var options []mal.UpdateMyAnimeListStatusOption
	options = append(options, mal.NumEpisodesWatched(am.show.ep))

	if ml.status == mal.AnimeStatusCompleted {
		options = append(options, mal.IsRewatching(true))
	}

	ml.status = mal.AnimeStatusWatching

	if ml.epNum == am.show.ep {
		ml.status = mal.AnimeStatusCompleted
		if ml.isRewatch {
			options = append(options, mal.IsRewatching(false))
			ml.rewatchNum++
			options = append(options, mal.NumTimesRewatched(ml.rewatchNum))
		}
	}

	options = append(options, ml.status)

	l, _, err := am.client.Anime.UpdateMyListStatus(ctx, am.malid, options...)
	if err != nil {
		return err
	}

	logUpdate(ml, l)
	return nil
}

func (am *AnimeUpdate) checkAnime(ctx context.Context) (*MyList, error) {

	a, _, err := am.client.Anime.Details(ctx, am.malid, mal.Fields{"num_episodes", "title", "my_list_status"})
	if err != nil {
		return nil, err
	}

	ml := &MyList{
		status:     a.MyListStatus.Status,
		isRewatch:  a.MyListStatus.IsRewatching,
		rewatchNum: a.MyListStatus.NumTimesRewatched,
		epNum:      a.NumEpisodes,
		title:      a.Title,
	}

	return ml, nil
}

func (am *AnimeUpdate) updateRating(ctx context.Context, r float32) error {

	ml, err := am.checkAnime(ctx)
	if err != nil {
		return err
	}

	l, _, err := am.client.Anime.UpdateMyListStatus(ctx, am.malid, mal.Score(r))
	if err != nil {
		return err
	}

	logUpdate(ml, l)
	return nil
}

func (am *AnimeUpdate) getStartID(ctx context.Context, multi bool) {

	for _, anime := range am.animeSeasons.Seasons {
		if am.show.season == anime.Season {
			if multi {
				am.malid = anime.MalID
				am.start = anime.Start
			} else {
				if am.show.ep >= anime.Start {
					am.malid = anime.MalID
					am.start = anime.Start
				}
			}
		}
	}

	am.start = updateStart(ctx, am.start)
}

func synonymExists(ctx context.Context, s []string, title string) bool {

	for _, v := range s {
		if v == title {
			return true
		}
	}
	return false
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
