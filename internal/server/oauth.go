package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type MalCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func NewOauth2Client(ctx context.Context) *http.Client {

	creds := &MalCredentials{}
	err := jsonUnmarshal("./credentials.json", creds)
	if err != nil {
		log.Fatalf("Not able to unmarshal mal credentials: %v", err)
	}

	cfg := &oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	t := &oauth2.Token{}
	err = jsonUnmarshal("./token.json", t)
	if err != nil {
		log.Println("Token not stored! Visit below URL to get token: ")
		t = getToken(ctx, cfg)
		saveToken(t)
	}

	fresh_token, err := cfg.TokenSource(ctx, t).Token()
	if err == nil && (fresh_token != t) {
		saveToken(fresh_token)
	}

	client := cfg.Client(ctx, fresh_token)

	return client
}

func getToken(ctx context.Context, cfg *oauth2.Config) *oauth2.Token {

	var (
		pkce          string                = randomString(128)
		state         string                = randomString(32)
		CodeChallenge oauth2.AuthCodeOption = oauth2.SetAuthURLParam("code_challenge", pkce)
		ResponseType  oauth2.AuthCodeOption = oauth2.SetAuthURLParam("response_type", "code")
		GrantType     oauth2.AuthCodeOption = oauth2.SetAuthURLParam("grant_type", "authorization_code")
		CodeVerify    oauth2.AuthCodeOption = oauth2.SetAuthURLParam("code_verifier", pkce)
	)

	url := cfg.AuthCodeURL(state, CodeChallenge, ResponseType)

	fmt.Println(url)

	fmt.Println("Enter code from MAL auth server below:")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	code := sc.Text()
	if sc.Err() != nil {
		log.Fatalf("Could not read code! %v", sc.Err())
	}

	token, err := cfg.Exchange(ctx, code, GrantType, CodeVerify)
	if err != nil {
		log.Fatalln("Could not get access token!", err)
	}

	return token

}
