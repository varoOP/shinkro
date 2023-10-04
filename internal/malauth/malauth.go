package malauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/varoOP/shinkro/internal/database"
	"golang.org/x/oauth2"
)

func NewOauth2Client(ctx context.Context, db *database.DB) (*http.Client, error) {
	creds, err := db.GetMalCreds(ctx)
	if err != nil {
		return nil, err
	}

	cfg := GetCfg(creds)
	t := &oauth2.Token{}
	err = json.Unmarshal([]byte(creds["access_token"]), t)
	if err != nil {
		return nil, err
	}

	fresh_token, err := cfg.TokenSource(ctx, t).Token()
	if err != nil {
		return nil, err
	}

	if err == nil && (fresh_token != t) {
		SaveToken(fresh_token, creds["client_id"], creds["client_secret"], db)
	}

	client := cfg.Client(ctx, fresh_token)
	return client, nil
}

func GetOauth(ctx context.Context, client_id, client_secret string) (*oauth2.Config, map[string]string) {
	var (
		pkce          string                = randomString(128)
		state         string                = randomString(32)
		CodeChallenge oauth2.AuthCodeOption = oauth2.SetAuthURLParam("code_challenge", pkce)
		ResponseType  oauth2.AuthCodeOption = oauth2.SetAuthURLParam("response_type", "code")
		creds                               = map[string]string{
			"client_id":     client_id,
			"client_secret": client_secret,
		}
	)

	cfg := GetCfg(creds)
	return cfg, map[string]string{
		"AuthCodeURL": cfg.AuthCodeURL(state, CodeChallenge, ResponseType),
		"state":       state,
		"pkce":        pkce,
	}
}

func SaveToken(token *oauth2.Token, client_id, client_secret string, db *database.DB) {
	t, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}

	db.UpdateMalAuth(map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"access_token":  string(t),
	})
}

func GetCfg(creds map[string]string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     creds["client_id"],
		ClientSecret: creds["client_secret"],
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

func randomString(l int) string {
	random := make([]byte, l)
	_, err := rand.Read(random)
	if err != nil {
		log.Fatalln(err)
	}

	return base64.URLEncoding.EncodeToString(random)[:l]
}
