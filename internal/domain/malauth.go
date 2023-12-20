package domain

import (
	"context"

	"golang.org/x/oauth2"
)

type MalAuthRepo interface {
	Store(ctx context.Context, ma *MalAuth) error
	Get(ctx context.Context) (*MalAuth, error)
}

type MalAuth struct {
	Id          int
	Config      oauth2.Config
	AccessToken oauth2.Token
}

type MalAuthOpts struct {
	MalAuth     *MalAuth
	Verifier    string
	State       string
	AuthCodeUrl string
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
