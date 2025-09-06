package domain

import "context"

type PlexSettingsRepo interface {
	Store(ctx context.Context, userID int, ps PlexSettings) (*PlexSettings, error)
	Update(ctx context.Context, userID int, ps PlexSettings) (*PlexSettings, error)
	Get(ctx context.Context, userID int) (*PlexSettings, error)
	Delete(ctx context.Context, userID int) error
}

type PlexSettings struct {
	UserID            int      `json:"user_id"`
	Host              string   `json:"host"`
	Port              int      `json:"port"`
	TLS               bool     `json:"tls"`
	TLSSkip           bool     `json:"tls_skip"`
	AnimeLibraries    []string `json:"anime_libs"`
	PlexUser          string   `json:"plex_user"`
	PlexClientEnabled bool     `json:"plex_client_enabled"`
	Token             []byte   `json:"-"`
	TokenIV           []byte   `json:"-"`
	ClientID          string   `json:"client_id"`
}

func NewPlexSettings(userID int, host, plexUser, clientID string, token, tokenIV []byte, port int, animeLibs []string, pce, tls, tlsSkip bool) *PlexSettings {
	return &PlexSettings{
		UserID:            userID,
		Host:              host,
		Port:              port,
		TLS:               tls,
		TLSSkip:           tlsSkip,
		AnimeLibraries:    animeLibs,
		PlexUser:          plexUser,
		PlexClientEnabled: pce,
		Token:             token,
		TokenIV:           tokenIV,
		ClientID:          clientID,
	}
}
