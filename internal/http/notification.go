package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/varoOP/shinkro/internal/domain"

	"github.com/go-chi/chi/v5"
)

type notificationService interface {
	Find(context.Context, domain.NotificationQueryParams) ([]domain.Notification, int, error)
	FindByID(ctx context.Context, userID, id int) (*domain.Notification, error)
	Store(ctx context.Context, userID int, notification *domain.Notification) error
	Update(ctx context.Context, userID int, notification *domain.Notification) error
	Delete(ctx context.Context, userID, id int) error
	Test(ctx context.Context, notification *domain.Notification) error
}

type notificationHandler struct {
	encoder encoder
	service notificationService
}

func newNotificationHandler(encoder encoder, service notificationService) *notificationHandler {
	return &notificationHandler{
		encoder: encoder,
		service: service,
	}
}

func (h notificationHandler) Routes(r chi.Router) {
	r.Get("/", h.list)
	r.Post("/", h.store)
	r.Post("/test", h.test)

	r.Route("/{notificationID}", func(r chi.Router) {
		r.Get("/", h.findByID)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})
}

func (h notificationHandler) list(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    "UNAUTHORIZED",
			"message": "User not authenticated",
		})
		return
	}

	list, _, err := h.service.Find(r.Context(), domain.NotificationQueryParams{UserID: userID})
	if err != nil {
		h.encoder.StatusNotFound(w)
		return
	}

	h.encoder.StatusResponse(w, http.StatusOK, list)
}

func (h notificationHandler) store(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    "UNAUTHORIZED",
			"message": "User not authenticated",
		})
		return
	}

	var data *domain.Notification
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.encoder.Error(w, err)
		return
	}

	err = h.service.Store(r.Context(), userID, data)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponse(w, http.StatusCreated, data)
}

func (h notificationHandler) findByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    "UNAUTHORIZED",
			"message": "User not authenticated",
		})
		return
	}

	notificationID, err := strconv.Atoi(chi.URLParam(r, "notificationID"))
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	notification, err := h.service.FindByID(r.Context(), userID, notificationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.encoder.NotFoundErr(w, errors.New(fmt.Sprintf("notification with id %d not found", notificationID)))
			return
		}

		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponse(w, http.StatusNoContent, notification)
}

func (h notificationHandler) update(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    "UNAUTHORIZED",
			"message": "User not authenticated",
		})
		return
	}

	var data *domain.Notification
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.encoder.Error(w, err)
		return
	}

	err = h.service.Update(r.Context(), userID, data)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponse(w, http.StatusOK, data)
}

func (h notificationHandler) delete(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    "UNAUTHORIZED",
			"message": "User not authenticated",
		})
		return
	}

	notificationID, err := strconv.Atoi(chi.URLParam(r, "notificationID"))
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), userID, notificationID); err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponse(w, http.StatusNoContent, nil)
}

func (h notificationHandler) test(w http.ResponseWriter, r *http.Request) {
	var data *domain.Notification
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.encoder.Error(w, err)
		return
	}

	if err := h.service.Test(r.Context(), data); err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.NoContent(w)
}
