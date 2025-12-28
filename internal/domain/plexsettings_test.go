package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlexSettings(t *testing.T) {
	tests := []struct {
		name       string
		host       string
		plexUser   string
		clientID   string
		token      []byte
		tokenIV    []byte
		port       int
		animeLibs  []string
		pce        bool
		tls        bool
		tlsSkip    bool
		validate   func(*testing.T, *PlexSettings)
	}{
		{
			name:     "creates plex settings with all fields",
			host:     "localhost",
			plexUser: "TestUser",
			clientID: "test-client-id",
			token:    []byte("test-token"),
			tokenIV:  []byte("test-token-iv"),
			port:     32400,
			animeLibs: []string{"Anime", "Anime Movies"},
			pce:       true,
			tls:       true,
			tlsSkip:   false,
			validate: func(t *testing.T, ps *PlexSettings) {
				assert.Equal(t, "localhost", ps.Host)
				assert.Equal(t, "TestUser", ps.PlexUser)
				assert.Equal(t, "test-client-id", ps.ClientID)
				assert.Equal(t, []byte("test-token"), ps.Token)
				assert.Equal(t, []byte("test-token-iv"), ps.TokenIV)
				assert.Equal(t, 32400, ps.Port)
				assert.Equal(t, []string{"Anime", "Anime Movies"}, ps.AnimeLibraries)
				assert.True(t, ps.PlexClientEnabled)
				assert.True(t, ps.TLS)
				assert.False(t, ps.TLSSkip)
			},
		},
		{
			name:     "creates plex settings with minimal fields",
			host:     "192.168.1.100",
			plexUser: "User",
			clientID: "client-123",
			token:    []byte{},
			tokenIV:  []byte{},
			port:     32400,
			animeLibs: []string{},
			pce:       false,
			tls:       false,
			tlsSkip:   false,
			validate: func(t *testing.T, ps *PlexSettings) {
				assert.Equal(t, "192.168.1.100", ps.Host)
				assert.Equal(t, "User", ps.PlexUser)
				assert.Equal(t, "client-123", ps.ClientID)
				assert.Equal(t, []byte{}, ps.Token)
				assert.Equal(t, []byte{}, ps.TokenIV)
				assert.Equal(t, 32400, ps.Port)
				assert.Equal(t, []string{}, ps.AnimeLibraries)
				assert.False(t, ps.PlexClientEnabled)
				assert.False(t, ps.TLS)
				assert.False(t, ps.TLSSkip)
			},
		},
		{
			name:     "creates plex settings with TLS skip enabled",
			host:     "plex.example.com",
			plexUser: "Admin",
			clientID: "admin-client",
			token:    []byte("admin-token"),
			tokenIV:  []byte("admin-iv"),
			port:     443,
			animeLibs: []string{"Anime"},
			pce:       true,
			tls:       true,
			tlsSkip:   true,
			validate: func(t *testing.T, ps *PlexSettings) {
				assert.Equal(t, "plex.example.com", ps.Host)
				assert.Equal(t, 443, ps.Port)
				assert.True(t, ps.TLS)
				assert.True(t, ps.TLSSkip)
			},
		},
		{
			name:     "creates plex settings with empty strings",
			host:     "",
			plexUser: "",
			clientID: "",
			token:    nil,
			tokenIV:  nil,
			port:     0,
			animeLibs: nil,
			pce:       false,
			tls:       false,
			tlsSkip:   false,
			validate: func(t *testing.T, ps *PlexSettings) {
				assert.Equal(t, "", ps.Host)
				assert.Equal(t, "", ps.PlexUser)
				assert.Equal(t, "", ps.ClientID)
				assert.Nil(t, ps.Token)
				assert.Nil(t, ps.TokenIV)
				assert.Equal(t, 0, ps.Port)
				assert.Nil(t, ps.AnimeLibraries)
			},
		},
		{
			name:     "creates plex settings with single anime library",
			host:     "localhost",
			plexUser: "User",
			clientID: "client",
			token:    []byte("token"),
			tokenIV:  []byte("iv"),
			port:     32400,
			animeLibs: []string{"Anime"},
			pce:       false,
			tls:       false,
			tlsSkip:   false,
			validate: func(t *testing.T, ps *PlexSettings) {
				assert.Equal(t, []string{"Anime"}, ps.AnimeLibraries)
				assert.Equal(t, 1, len(ps.AnimeLibraries))
			},
		},
		{
			name:     "creates plex settings with multiple anime libraries",
			host:     "localhost",
			plexUser: "User",
			clientID: "client",
			token:    []byte("token"),
			tokenIV:  []byte("iv"),
			port:     32400,
			animeLibs: []string{"Anime", "Anime Movies", "Anime TV"},
			pce:       false,
			tls:       false,
			tlsSkip:   false,
			validate: func(t *testing.T, ps *PlexSettings) {
				assert.Equal(t, []string{"Anime", "Anime Movies", "Anime TV"}, ps.AnimeLibraries)
				assert.Equal(t, 3, len(ps.AnimeLibraries))
			},
		},
		{
			name:     "creates plex settings with custom port",
			host:     "localhost",
			plexUser: "User",
			clientID: "client",
			token:    []byte("token"),
			tokenIV:  []byte("iv"),
			port:     8080,
			animeLibs: []string{},
			pce:       false,
			tls:       false,
			tlsSkip:   false,
			validate: func(t *testing.T, ps *PlexSettings) {
				assert.Equal(t, 8080, ps.Port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPlexSettings(
				tt.host,
				tt.plexUser,
				tt.clientID,
				tt.token,
				tt.tokenIV,
				tt.port,
				tt.animeLibs,
				tt.pce,
				tt.tls,
				tt.tlsSkip,
			)
			require.NotNil(t, result)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

