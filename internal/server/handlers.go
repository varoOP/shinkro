package server

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/domain"
	"github.com/varoOP/shinkuro/pkg/plex"
)

func Plex(db *database.DB, cfg *domain.Config, log *zerolog.Logger, n *domain.Notification) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		a := domain.NewAnimeUpdate(db, cfg, log, n)
		var err error
		p := r.Context().Value(PlexPayload).(*plex.PlexWebhook)
		a.Event = p.Event
		a.Rating = p.Rating
		a.Mapping, err = domain.NewAnimeSeasonMap(a.Config)
		if err != nil {
			a.Log.Error().Err(err).Msg("unable to load custom mapping")
			return
		}

		a.InMap, a.Anime = a.Mapping.CheckAnimeMap(p.Metadata.GrandparentTitle)
		a.Media, err = database.NewMedia(p.Metadata.GUID.GUID, p.Metadata.Type, r.Context().Value(MediaTitle).(string))
		if err != nil {
			a.Log.Error().Err(err).Msg("unable to parse media")
			return
		}

		err = a.SendUpdate(r.Context())
		if err.Error() == "complete" {
			return
		}

		notify(&a, err)
		if err != nil {
			a.Log.Error().Err(err).Msg("failed to send update to myanimelist")
			return
		}

		a.Log.Info().
			Str("status", string(a.Malresp.Status)).
			Int("score", a.Malresp.Score).
			Int("episdoesWatched", a.Malresp.NumEpisodesWatched).
			Int("timesRewatched", a.Malresp.NumTimesRewatched).
			Str("startDate", a.Malresp.StartDate).
			Str("finishDate", a.Malresp.FinishDate).
			Msg("updated myanimelist successfully")

		w.WriteHeader(http.StatusNoContent)
	}
}
