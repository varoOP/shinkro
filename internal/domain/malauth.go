package domain

import (
	"context"

	"golang.org/x/oauth2"
)

type MalAuthRepo interface {
	Store(ctx context.Context, ma *MalAuth) error
	Get(ctx context.Context) (*MalAuth, error)
	Delete(ctx context.Context) error
}

type MalAuth struct {
	Id          int
	Config      oauth2.Config
	AccessToken []byte
	TokenIV     []byte
}

type MalAuthURLs string

const AuthURL MalAuthURLs = "https://myanimelist.net/v1/oauth2/authorize"
const TokenURL MalAuthURLs = "https://myanimelist.net/v1/oauth2/token"

func NewMalAuth(clientID, clientSecret string, accessToken, tokenIV []byte) *MalAuth {
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
		AccessToken: accessToken,
		TokenIV:     tokenIV,
	}
}
