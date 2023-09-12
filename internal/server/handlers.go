package server

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/pkg/plex"
)

func Plex(db *database.DB, cfg *domain.Config, log *zerolog.Logger, n *domain.Notification) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		a := domain.NewAnimeUpdate(db, cfg, log, n)
		a.Plex = r.Context().Value(domain.PlexPayload).(*plex.PlexWebhook)
		err = a.SendUpdate(r.Context())
		if err != nil && err.Error() == "complete" {
			return
		}

		notify(&a, err)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			a.Log.Error().Stack().Err(err).Msg("failed to send update to myanimelist")
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
