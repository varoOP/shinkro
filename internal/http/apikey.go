package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/varoOP/shinkro/internal/domain"
)

type apikeyService interface {
	List(ctx context.Context) ([]domain.APIKey, error)
	Store(ctx context.Context, key *domain.APIKey) error
	Delete(ctx context.Context, key string) error
	ValidateAPIKey(ctx context.Context, token string) bool
	GetUserIDByAPIKey(ctx context.Context, token string) (int, error)
}

type apikeyHandler struct {
	encoder encoder
	service apikeyService
}

func newAPIKeyHandler(encoder encoder, service apikeyService) *apikeyHandler {
	return &apikeyHandler{
		encoder: encoder,
		service: service,
	}
}

func (h apikeyHandler) Routes(r chi.Router) {
	r.Get("/", h.list)
	r.Post("/", h.store)
	r.Delete("/{apikey}", h.delete)
}

func (h apikeyHandler) list(w http.ResponseWriter, r *http.Request) {
	keys, err := h.service.List(r.Context())
	if err != nil {
		h.encoder.StatusError(w, http.StatusInternalServerError, err)
		return
	}

	render.JSON(w, r, keys)
}

func (h apikeyHandler) store(w http.ResponseWriter, r *http.Request) {
	var data domain.APIKey
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.encoder.Error(w, err)
		return
	}

	if err := h.service.Store(r.Context(), &data); err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponse(w, http.StatusCreated, data)
}

func (h apikeyHandler) delete(w http.ResponseWriter, r *http.Request) {
	apiKey := chi.URLParam(r, "apikey")

	if err := h.service.Delete(r.Context(), apiKey); err != nil {
		if errors.Is(err, errors.New("record not found")) {
			h.encoder.NotFoundErr(w, errors.Errorf("api key %s not found", apiKey))
			return
		}

		h.encoder.Error(w, err)
		return
	}

	h.encoder.NoContent(w)
}
