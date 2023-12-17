package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/varoOP/shinkro/internal/domain"
)

type plexService interface {
	Store(ctx context.Context, plex *domain.Plex) error
	Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error)
}

type plexHandler struct {
	encoder encoder
	service plexService
}

func newPlexHandler(encoder encoder, service plexService) *plexHandler {
	return &plexHandler{
		encoder: encoder,
		service: service,
	}
}

func (h plexHandler) Routes(r chi.Router) {
	r.Get("/", h.GetPlex)
	r.With(middleware.AllowContentType("application/json", "multipart/form-data"), parsePlexPayload).Post("/", h.PostPlex)
}

func (h plexHandler) GetPlex(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idP)
	if err != nil && idP != "" {
		h.encoder.StatusResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code":    "BAD_REQUEST_PARAMS",
			"message": "id parameter is invalid",
		})
		return
	}

	plex, err := h.service.Get(r.Context(), &domain.GetPlexRequest{Id: id})
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	ret := struct {
		Data *domain.Plex `json:"data"`
	}{
		Data: plex,
	}

	h.encoder.StatusResponse(w, http.StatusOK, ret)

}

func (h plexHandler) PostPlex(w http.ResponseWriter, r *http.Request) {
	plex := r.Context().Value(domain.PlexPayload).(*domain.Plex)
	err := h.service.Store(r.Context(), plex)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code":    "BAD_REQUEST_PARAMS",
			"message": err.Error(),
		})
		return
	}

	h.encoder.StatusCreated(w)
}
