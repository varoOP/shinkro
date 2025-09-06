package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/domain"
)

type animeupdateService interface {
	Count(ctx context.Context) (int, error)
	GetRecentUnique(ctx context.Context, userID int, limit int) ([]*domain.AnimeUpdate, error)
	GetByPlexID(ctx context.Context, plexID int64) (*domain.AnimeUpdate, error)
}

type GetRecentAnimeResponse struct {
	AnimeUpdates []domain.RecentAnimeItem `json:"animeUpdates"`
}

type animeupdateHandler struct {
	encoder encoder
	service animeupdateService
}

func newAnimeupdateHandler(encoder encoder, service animeupdateService) *animeupdateHandler {
	return &animeupdateHandler{
		encoder: encoder,
		service: service,
	}
}

func (h animeupdateHandler) Routes(r chi.Router) {
	r.Get("/count", h.getCount)
	r.Get("/recent", h.getRecent)
	r.Get("/byPlexId", h.getByPlexID)
}

func (h animeupdateHandler) getCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.Count(r.Context())
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}
	h.encoder.StatusResponse(w, http.StatusOK, map[string]interface{}{
		"count": count,
	})
}

func (h animeupdateHandler) getRecent(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)
	
	// Parse limit with default value
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		log.Error().Err(err).Msg("error getting user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	animeUpdates, err := h.service.GetRecentUnique(r.Context(), userID, limit)
	if err != nil {
		log.Error().Err(err).Msg("error getting recent unique anime updates")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Transform domain.AnimeUpdate to RecentAnimeItem
	items := make([]domain.RecentAnimeItem, 0, len(animeUpdates))
	for _, update := range animeUpdates {
		items = append(items, domain.RecentAnimeItem{
			AnimeStatus:     string(update.ListDetails.Status),
			FinishDate:      update.ListStatus.FinishDate,
			LastUpdated:     update.Timestamp.Format("2006-01-02T15:04:05Z"),
			MalId:           update.MALId,
			PictureUrl:      update.ListDetails.PictureURL,
			Rating:          update.ListStatus.Score,
			RewatchNum:      update.ListDetails.RewatchNum,
			StartDate:       update.ListStatus.StartDate,
			Title:           update.ListDetails.Title,
			TotalEpisodeNum: update.ListDetails.TotalEpisodeNum,
			WatchedNum:      update.ListDetails.WatchedNum,
		})
	}

	h.encoder.StatusResponse(w, http.StatusOK, GetRecentAnimeResponse{
		AnimeUpdates: items,
	})
}

func (h animeupdateHandler) getByPlexID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.encoder.StatusResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "missing id param"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "invalid id param"})
		return
	}
	update, err := h.service.GetByPlexID(r.Context(), id)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusNotFound, map[string]interface{}{"error": err.Error()})
		return
	}
	h.encoder.StatusResponse(w, http.StatusOK, update)
}
