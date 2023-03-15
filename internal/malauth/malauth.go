package malauth

import (
	"bufio"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/varoOP/shinkuro/internal/database"
	"golang.org/x/oauth2"
)

func NewOauth2Client(ctx context.Context, db *sql.DB) *http.Client {

	m := database.GetMalCreds(db)

	cfg := &oauth2.Config{
		ClientID:     m["client_id"],
		ClientSecret: m["client_secret"],
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	t := &oauth2.Token{}
	err := json.Unmarshal([]byte(m["access_token"]), t)
	if err != nil {
		log.Fatalln(err)
	}

	fresh_token, err := cfg.TokenSource(ctx, t).Token()
	if err == nil && (fresh_token != t) {
		saveToken(fresh_token, m["client_id"], m["client_secret"], db)
	}

	client := cfg.Client(ctx, fresh_token)

	return client
}

func NewMalAuth(db *sql.DB) {
	var (
		client_id     string
		client_secret string
	)

	fmt.Println("Enter MAL API client-id:")
	fmt.Scanln(&client_id)
	fmt.Println("Enter MAL API client-secret:")
	fmt.Scanln(&client_secret)

	if client_id == "" || client_secret == "" {
		log.Fatalf("client-id or client-secret not provided.")
	}

	t := getToken(context.Background(), client_id, client_secret)
	saveToken(t, client_id, client_secret, db)
}

func getToken(ctx context.Context, client_id, client_secret string) *oauth2.Token {

	var (
		pkce          string                = randomString(128)
		state         string                = randomString(32)
		CodeChallenge oauth2.AuthCodeOption = oauth2.SetAuthURLParam("code_challenge", pkce)
		ResponseType  oauth2.AuthCodeOption = oauth2.SetAuthURLParam("response_type", "code")
		GrantType     oauth2.AuthCodeOption = oauth2.SetAuthURLParam("grant_type", "authorization_code")
		CodeVerify    oauth2.AuthCodeOption = oauth2.SetAuthURLParam("code_verifier", pkce)
	)

	cfg := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	url := cfg.AuthCodeURL(state, CodeChallenge, ResponseType)

	fmt.Println(url)

	fmt.Println("Enter code from MAL auth server below:")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	code := sc.Text()
	if sc.Err() != nil {
		log.Fatalf("Could not read code! %v\n", sc.Err())
	}

	token, err := cfg.Exchange(ctx, code, GrantType, CodeVerify)
	if err != nil {
		log.Fatalln("Could not get access token!", err)
	}

	return token
}

func saveToken(token *oauth2.Token, client_id, client_secret string, db *sql.DB) {
	t, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}

	m := map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"access_token":  string(t),
	}

	database.UpdateMalAuth(m, db)
}

func randomString(l int) string {
	random := make([]byte, l)
	_, err := rand.Read(random)
	if err != nil {
		log.Fatalln(err)
	}

	return base64.URLEncoding.EncodeToString(random)[:l]
}
