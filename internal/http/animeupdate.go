package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/varoOP/shinkro/internal/domain"
)

type animeupdateService interface {
	Count(ctx context.Context) (int, error)
	GetRecentUnique(ctx context.Context, limit int) ([]*domain.AnimeUpdate, error)
	GetByPlexID(ctx context.Context, plexID int64) (*domain.AnimeUpdate, error)
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
	limit := 5
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	updates, err := h.service.GetRecentUnique(r.Context(), limit)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}
	// Only return the fields needed for the dashboard
	result := make([]map[string]interface{}, 0, len(updates))
	for _, u := range updates {
		result = append(result, map[string]interface{}{
			"malId":           u.MALId,
			"title":           u.ListDetails.Title,
			"pictureUrl":      u.ListDetails.PictureURL,
			"watchedNum":      u.ListStatus.NumEpisodesWatched,
			"totalEpisodeNum": u.ListDetails.TotalEpisodeNum,
			"lastUpdated":     u.ListStatus.UpdatedAt,
			"rating":          u.ListStatus.Score,
			"animeStatus":     u.ListStatus.Status,
			"startDate":       u.ListStatus.StartDate,
			"finishDate":      u.ListStatus.FinishDate,
			"rewatchNum":      u.ListDetails.RewatchNum,
		})
	}
	h.encoder.StatusResponse(w, http.StatusOK, result)
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
