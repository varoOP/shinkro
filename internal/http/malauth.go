package http

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strings"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkro/internal/domain"
	"golang.org/x/oauth2"
)

var (
	verifier    string
	state       string
	oauthConfig oauth2.Config
)

type malauthService interface {
	Store(ctx context.Context, ma *domain.MalAuth) error
	GetMalClient(ctx context.Context) (*mal.Client, error)
	NewMalAuthClient(ctx context.Context, clientId, clientSecret string) (*domain.MalAuthOpts, error)
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
	r.Get("/", h.malAuth)
	r.Post("/login", h.login)
	r.Get("/status", h.status)
	r.Get("/callback", h.callback)
}

func (h malauthHandler) login(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("clientID")
	clientSecret := r.FormValue("clientSecret")
	maopts, err := h.service.NewMalAuthClient(r.Context(), clientID, clientSecret)
	if err != nil {
		h.encoder.StatusResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	verifier = maopts.Verifier
	state = maopts.State
	oauthConfig = maopts.MalAuth.Config

	http.Redirect(w, r, maopts.AuthCodeUrl, http.StatusSeeOther)
}

func (h malauthHandler) callback(w http.ResponseWriter, r *http.Request) {
	baseURL := url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.Host,
	}

	newURL := baseURL.ResolveReference(&url.URL{Path: "/malauth/status"})

	code := r.URL.Query().Get("code")
	newState := r.URL.Query().Get("state")

	if code == "" || newState == "" || newState != state {
		http.Redirect(w, r, newURL.String(), http.StatusSeeOther)
		return
	}

	grantType := oauth2.SetAuthURLParam("grant_type", "authorization_code")
	codeVerify := oauth2.SetAuthURLParam("code_verifier", verifier)
	token, err := oauthConfig.Exchange(r.Context(), code, grantType, codeVerify)
	if err != nil {
		http.Redirect(w, r, newURL.String(), http.StatusSeeOther)
		return
	}

	h.service.Store(r.Context(), &domain.MalAuth{
		Id:          1,
		Config:      oauthConfig,
		AccessToken: *token,
	})

	http.Redirect(w, r, newURL.String(), http.StatusSeeOther)
}

func (h malauthHandler) status(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("malauthstatus").Parse(malauth_statustpl)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	isAuthenticated := false
	c, err := h.service.GetMalClient(r.Context())
	if err != nil {
		return
	}

	_, _, err = c.User.MyInfo(r.Context())
	if err == nil {
		isAuthenticated = true
	}

	segments := strings.Split(r.URL.Path, "/")
	if len(segments) > 1 {
		segments = segments[:len(segments)-1]
	}
	newPath := path.Join(segments...)
	newURL := *r.URL
	newURL.Path = newPath

	data := authPageData{
		IsAuthenticated: isAuthenticated,
		RetryURL:        newPath,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h malauthHandler) malAuth(w http.ResponseWriter, r *http.Request) {

	c, err := h.service.GetMalClient(r.Context())
	if err == nil {
		_, _, err := c.User.MyInfo(r.Context())
		if err == nil {
			w.Write([]byte("Authentication with myanimelist is successful."))
			return
		}
	}

	baseURL := url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.Host,
	}

	newURL := baseURL.ResolveReference(&url.URL{Path: "/malauth/login"})
	data := authPageData{
		ActionURL: newURL.String(),
	}

	tmpl, err := template.New("malauth").Parse(malauthtpl)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
