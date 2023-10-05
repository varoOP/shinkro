package server

import (
	"errors"
	"net/http"
	"text/template"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/pkg/plex"
	"golang.org/x/oauth2"
)

var (
	pkce          string
	state         string
	authConfig    *oauth2.Config
	CodeChallenge oauth2.AuthCodeOption
	ResponseType  oauth2.AuthCodeOption = oauth2.SetAuthURLParam("response_type", "code")
	GrantType     oauth2.AuthCodeOption = oauth2.SetAuthURLParam("grant_type", "authorization_code")
	CodeVerify    oauth2.AuthCodeOption
)

type AuthPageData struct {
	IsAuthenticated bool
	ActionURL       string
	RetryURL        string
}

func Plex(db *database.DB, cfg *domain.Config, log *zerolog.Logger, n *domain.Notification) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		a := domain.NewAnimeUpdate(db, cfg, log, n)
		a.Plex = r.Context().Value(domain.PlexPayload).(*plex.PlexWebhook)
		err = a.SendUpdate(r.Context())
		if err != nil && err.Error() == "complete" {
			return
		}

		notify(&a, err)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			a.Log.Error().Stack().Err(err).Msg("failed to send update to myanimelist")
			return
		}

		a.Log.Info().
			Str("title", string(a.Media.Title)).
			Interface("listStatus", a.Malresp).
			Msg("Updated myanimelist successfully!")

		w.WriteHeader(http.StatusNoContent)
	}
}

func MalAuthLogin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var authMap map[string]string
			clientID := r.FormValue("clientID")
			clientSecret := r.FormValue("clientSecret")
			authConfig, authMap = malauth.GetOauth(r.Context(), clientID, clientSecret)
			pkce = authMap["pkce"]
			state = authMap["state"]
			http.Redirect(w, r, authMap["AuthCodeURL"], http.StatusFound)
			return
		}
	}
}

func MalAuthCallback(cfg *domain.Config, db *database.DB, log *zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var code string
		u := joinUrlPath(cfg.BaseUrl, "/malauth/status")
		q := r.URL.Query()
		if len(q["code"]) >= 1 && len(q["state"]) >= 1 {
			code = q["code"][0]
			if state != q["state"][0] {
				http.Redirect(w, r, u, http.StatusSeeOther)
				log.Error().Err(errors.New("state did not match")).Str("state", q["state"][0]).Msg("")
				return
			}
		}

		CodeVerify = oauth2.SetAuthURLParam("code_verifier", pkce)
		token, err := authConfig.Exchange(r.Context(), code, GrantType, CodeVerify)
		if err != nil {
			http.Redirect(w, r, u, http.StatusSeeOther)
			log.Error().Err(err).Msg("")
			return
		}

		malauth.SaveToken(token, authConfig.ClientID, authConfig.ClientSecret, db)
		http.Redirect(w, r, u, http.StatusSeeOther)
	}
}

func MalAuthStatus(cfg *domain.Config, db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tmpl, err := template.New("malauthstatus").Parse(malauth_statustpl)
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}

		isAuthenticated := false
		client, _ := malauth.NewOauth2Client(r.Context(), db)
		c := mal.NewClient(client)
		_, _, err = c.User.MyInfo(r.Context())
		if err == nil {
			isAuthenticated = true
		}

		data := AuthPageData{
			IsAuthenticated: isAuthenticated,
			RetryURL:        joinUrlPath(cfg.BaseUrl, "/malauth"),
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}

func MalAuth(cfg *domain.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := AuthPageData{
			ActionURL: joinUrlPath(cfg.BaseUrl, "/malauth/login"),
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
}

func NotFound(cfg *domain.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := joinUrlPath(cfg.BaseUrl, "/malauth")
		http.Redirect(w, r, u, http.StatusSeeOther)
	}
}
