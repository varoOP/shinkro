package domain

import (
	"context"

	"golang.org/x/oauth2"
)

type MalAuthRepo interface {
	Store(ctx context.Context, ma *MalAuth) error
	Get(ctx context.Context) (*MalAuth, error)
	StoreMalAuthOpts(ctx context.Context, mo *MalAuthOpts) error
	GetMalAuthOpts(ctx context.Context) (*MalAuthOpts, error)
}

type MalAuth struct {
	Id          int
	Config      oauth2.Config
	AccessToken oauth2.Token
}

type MalAuthOpts struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	Verifier     string `json:"verifier"`
	State        string `json:"state"`
	Code         string `json:"code"`
}

type MalAuthURLs string

const AuthURL MalAuthURLs = "https://myanimelist.net/v1/oauth2/authorize"
const TokenURL MalAuthURLs = "https://myanimelist.net/v1/oauth2/token"

func NewMalAuth(clientID, clientSecret string) *MalAuth {
	return &MalAuth{
		Id: 1,
		Config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   string(AuthURL),
				TokenURL:  string(TokenURL),
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
	}
}

func NewMalAuthOpts(clientID, clientSecret, verifier, state, code string) *MalAuthOpts {
	return &MalAuthOpts{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Verifier:     verifier,
		State:        state,
		Code:         code,
	}
}
