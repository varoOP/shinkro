package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/varoOP/shinkuro/pkg/mal"
	"golang.org/x/oauth2"
)

var (
	Pkce         string         = randomString(128)
	State        string         = randomString(32)
	TokenCh                     = make(chan *oauth2.Token, 1)
	Oauth2Config *oauth2.Config = NewOauth2Config()
	Oauth2Client *http.Client
)

func NewOauth2Config() *oauth2.Config {

	creds := &mal.MalCredentials{}
	err := jsonUnmarshalHelp("./credentials.json", creds)
	if err != nil {
		log.Fatalf("Not able to unmarshal mal credentials: %v", err)
	}

	return &oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

}

func GetToken(cfg *oauth2.Config) *oauth2.Token {

	var (
		CodeChallenge oauth2.AuthCodeOption = oauth2.SetAuthURLParam("code_challenge", Pkce)
		ResponseType  oauth2.AuthCodeOption = oauth2.SetAuthURLParam("response_type", "code")
	)

	url := cfg.AuthCodeURL(State, CodeChallenge, ResponseType)

	fmt.Println(url)

	token := <-TokenCh

	return token

}

func NewOauth2Client(cfg *oauth2.Config) *http.Client {

	t := &oauth2.Token{}
	err := jsonUnmarshalHelp("./token.json", t)
	if err != nil {
		log.Println("Token not stored! Visit below URL to get token: ")
		t = GetToken(cfg)
		go saveToken(t)
	}

	client := cfg.Client(context.Background(), t)

	return client

}

func jsonUnmarshalHelp(path string, a interface{}) error {

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

func saveToken(token *oauth2.Token) {

	tokenJ, err := os.Create("./token.json")
	if err != nil {
		log.Fatalln(err)
	}

	defer tokenJ.Close()

	e := json.NewEncoder(tokenJ)
	e.SetIndent("", "  ")
	e.Encode(token)
	log.Println("Token saved.")

}
