package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type animeupdateService interface {
	Count(ctx context.Context) (int, error)
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
