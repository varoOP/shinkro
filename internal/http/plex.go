package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/domain"
)

type plexService interface {
	Store(ctx context.Context, plex *domain.Plex) error
	Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error)
	ProcessPlex(ctx context.Context, plex *domain.Plex, agent *domain.PlexSupportedAgents) error
	GetPlexSettings(ctx context.Context) (*domain.PlexSettings, error)
	CountScrobbleEvents(ctx context.Context) (int, error)
	CountRateEvents(ctx context.Context) (int, error)
	GetPlexHistory(ctx context.Context, req *domain.PlexHistoryRequest) (*domain.PlexHistoryResponse, error)
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
	r.Get("/", h.getPlex)
	r.Get("/count", h.getCounts)
	r.With(middleware.AllowContentType("application/json", "multipart/form-data"), parsePlexPayload).Post("/", h.postPlex)
	r.Get("/history", h.getHistory)
}

func (h plexHandler) getPlex(w http.ResponseWriter, r *http.Request) {
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

func (h plexHandler) postPlex(w http.ResponseWriter, r *http.Request) {
	plex := r.Context().Value(domain.PlexPayload).(*domain.Plex)
	plexSettings, err := h.service.GetPlexSettings(r.Context())
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code":    "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	agent, err := plex.CheckPlex(plexSettings)
	if err != nil {
		log := hlog.FromRequest(r)
		log.Debug().Err(err).Msg("Plex payload not sent for processing")
		h.encoder.StatusResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code":    "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	err = h.service.Store(r.Context(), plex)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusBadRequest, map[string]interface{}{
			"code":    "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	err = h.service.ProcessPlex(r.Context(), plex, &agent)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}
	// if !h.service.CheckPlex(plex) {
	// 	h.encoder.StatusResponse(w, http.StatusOK, map[string]interface{}{
	// 		"code":    "OK",
	// 		"message": "Check Plex false",
	// 	})
	// 	return
	// }

	h.encoder.StatusCreated(w)
}

func (h plexHandler) getCounts(w http.ResponseWriter, r *http.Request) {
	countScrobble, err := h.service.CountScrobbleEvents(r.Context())
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}
	countRate, err := h.service.CountRateEvents(r.Context())
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}
	h.encoder.StatusResponse(w, http.StatusOK, map[string]interface{}{
		"countScrobble": countScrobble,
		"countRate":     countRate,
	})
}

func (h plexHandler) getHistory(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	req := &domain.PlexHistoryRequest{
		Limit:    limit,
		Offset:   offset,
		Cursor:   r.URL.Query().Get("cursor"),
		Search:   r.URL.Query().Get("search"),
		Status:   r.URL.Query().Get("status"),
		Event:    r.URL.Query().Get("event"),
		FromDate: r.URL.Query().Get("from"),
		ToDate:   r.URL.Query().Get("to"),
		Type:     r.URL.Query().Get("type"),
	}

	// Default to "table" if not specified
	if req.Type == "" {
		req.Type = "table"
	}

	response, err := h.service.GetPlexHistory(r.Context(), req)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	h.encoder.StatusResponse(w, http.StatusOK, response)
}
