package server

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func NewOauth2Client(ctx context.Context, client_id string, client_secret string, token_path string) *http.Client {

	cfg := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	t := &oauth2.Token{}
	err := jsonUnmarshal(token_path, t)
	if err != nil {
		log.Println("Token not stored! Visit below URL to get token: ")
		t = getToken(ctx, cfg, token_path)
	}

	fresh_token, err := cfg.TokenSource(ctx, t).Token()
	if err == nil && (fresh_token != t) {
		saveToken(fresh_token, token_path)
	}

	client := cfg.Client(ctx, fresh_token)

	return client
}

func getToken(ctx context.Context, cfg *oauth2.Config, token_path string) *oauth2.Token {

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
		log.Fatalf("Could not read code! %v\n", sc.Err())
	}

	token, err := cfg.Exchange(ctx, code, GrantType, CodeVerify)
	if err != nil {
		log.Fatalln("Could not get access token!", err)
	}

	saveToken(token, token_path)

	return token

}

func jsonUnmarshal(path string, a interface{}) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, a)
	if err != nil {
		return err
	}

	return nil
}

func randomString(l int) string {
	random := make([]byte, l)
	_, err := rand.Read(random)
	if err != nil {
		log.Fatalln(err)
	}

	return base64.URLEncoding.EncodeToString(random)[:l]
}

func saveToken(token *oauth2.Token, token_path string) {

	tokenJ, err := os.Create(token_path)
	if err != nil {
		log.Fatalln(err)
	}

	defer tokenJ.Close()

	e := json.NewEncoder(tokenJ)
	e.SetIndent("", "  ")
	e.Encode(token)
	log.Println("Token saved at", token_path)

}
