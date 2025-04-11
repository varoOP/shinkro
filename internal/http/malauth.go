package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkro/internal/domain"
	"golang.org/x/oauth2"
	"net/http"
)

type malauthService interface {
	Store(ctx context.Context, ma *domain.MalAuth) error
	StoreMalAuthOpts(ctx context.Context, mo *domain.MalAuthOpts) error
	GetMalAuthOpts(ctx context.Context) (*domain.MalAuthOpts, error)
	GetMalClient(ctx context.Context) (*mal.Client, error)
}

type malauthHandler struct {
	encoder encoder
	service malauthService
}

func newmalauthHandler(encoder encoder, service malauthService) *malauthHandler {
	return &malauthHandler{
		encoder: encoder,
		service: service,
	}
}

func (h malauthHandler) Routes(r chi.Router) {
	r.Get("/test", h.test)
	r.Post("/opts", h.storeMalauthOpts)
	r.Get("/opts", h.getMalauthOpts)
	r.Post("/callback", h.callback)
}

func (h malauthHandler) storeMalauthOpts(w http.ResponseWriter, r *http.Request) {
	malauthopts := &domain.MalAuthOpts{}
	err := json.NewDecoder(r.Body).Decode(malauthopts)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	err = h.service.StoreMalAuthOpts(r.Context(), malauthopts)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponseMessage(w, http.StatusOK, "malauth opts stored")
}

func (h malauthHandler) getMalauthOpts(w http.ResponseWriter, r *http.Request) {
	malauthopts, err := h.service.GetMalAuthOpts(r.Context())
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusOK, map[string]interface{}{})
		return
	}

	h.encoder.StatusResponse(w, http.StatusOK, malauthopts)
}

func (h malauthHandler) callback(w http.ResponseWriter, r *http.Request) {
	malauthopts, err := h.service.GetMalAuthOpts(r.Context())
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&malauthopts)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	malauth := domain.NewMalAuth(malauthopts.ClientID, malauthopts.ClientSecret)
	grantType := oauth2.SetAuthURLParam("grant_type", "authorization_code")
	codeVerify := oauth2.SetAuthURLParam("code_verifier", malauthopts.Verifier)
	token, err := malauth.Config.Exchange(r.Context(), malauthopts.Code, grantType, codeVerify)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	malauth.AccessToken = *token
	err = h.service.Store(r.Context(), malauth)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponseMessage(w, http.StatusOK, "mal auth success")
}

func (h malauthHandler) test(w http.ResponseWriter, r *http.Request) {

	c, err := h.service.GetMalClient(r.Context())
	if err != nil {
		h.encoder.Error(w, err)
		return
	}
	_, _, err = c.User.MyInfo(r.Context())
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.StatusResponseMessage(w, http.StatusOK, "mal auth test success")

}
