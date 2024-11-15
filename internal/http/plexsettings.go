package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/varoOP/shinkro/internal/domain"
)

type plexsettingsService interface {
	Store(ctx context.Context, ps domain.PlexSettings) (*domain.PlexSettings, error)
	Get(ctx context.Context) (*domain.PlexSettings, error)
}

type plexsettingsHandler struct {
	encoder encoder
	service plexsettingsService
}

func newPlexsettingsHandler(encoder encoder, service plexsettingsService) *plexsettingsHandler {
	return &plexsettingsHandler{
		encoder: encoder,
		service: service,
	}
}

func (h plexsettingsHandler) Routes(r chi.Router) {
	r.Get("/", h.getPlexSettings)
	r.Post("/", h.postPlexSettings)
}

func (h plexsettingsHandler) getPlexSettings(w http.ResponseWriter, r *http.Request) {
	ps, err := h.service.Get(r.Context())
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	ret := struct {
		Data *domain.PlexSettings `json:"data"`
	}{
		Data: ps,
	}

	h.encoder.StatusResponse(w, http.StatusOK, ret)
}

func (h plexsettingsHandler) postPlexSettings(w http.ResponseWriter, r *http.Request) {
	var data domain.PlexSettings
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.encoder.Error(w, err)
		return
	}

	ps, err := h.service.Store(r.Context(), data)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponse(w, http.StatusCreated, ps)
}
