package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestNewMalAuth(t *testing.T) {
	tests := []struct {
		name          string
		clientID      string
		clientSecret  string
		accessToken   []byte
		tokenIV       []byte
		validate      func(*testing.T, *MalAuth)
	}{
		{
			name:         "creates mal auth with all fields",
			clientID:     "test-client-id",
			clientSecret: "test-client-secret",
			accessToken:  []byte("test-access-token"),
			tokenIV:      []byte("test-token-iv"),
			validate: func(t *testing.T, ma *MalAuth) {
				assert.Equal(t, 1, ma.Id)
				assert.Equal(t, "test-client-id", ma.Config.ClientID)
				assert.Equal(t, "test-client-secret", ma.Config.ClientSecret)
				assert.Equal(t, string(AuthURL), ma.Config.Endpoint.AuthURL)
				assert.Equal(t, string(TokenURL), ma.Config.Endpoint.TokenURL)
				assert.Equal(t, oauth2.AuthStyleInParams, ma.Config.Endpoint.AuthStyle)
				assert.Equal(t, []byte("test-access-token"), ma.AccessToken)
				assert.Equal(t, []byte("test-token-iv"), ma.TokenIV)
			},
		},
		{
			name:         "creates mal auth with empty strings",
			clientID:     "",
			clientSecret: "",
			accessToken:  []byte{},
			tokenIV:      []byte{},
			validate: func(t *testing.T, ma *MalAuth) {
				assert.Equal(t, 1, ma.Id)
				assert.Equal(t, "", ma.Config.ClientID)
				assert.Equal(t, "", ma.Config.ClientSecret)
				assert.Equal(t, []byte{}, ma.AccessToken)
				assert.Equal(t, []byte{}, ma.TokenIV)
				// OAuth endpoints should still be set
				assert.Equal(t, string(AuthURL), ma.Config.Endpoint.AuthURL)
				assert.Equal(t, string(TokenURL), ma.Config.Endpoint.TokenURL)
			},
		},
		{
			name:         "creates mal auth with nil tokens",
			clientID:     "client-id",
			clientSecret: "client-secret",
			accessToken:  nil,
			tokenIV:      nil,
			validate: func(t *testing.T, ma *MalAuth) {
				assert.Equal(t, 1, ma.Id)
				assert.Equal(t, "client-id", ma.Config.ClientID)
				assert.Equal(t, "client-secret", ma.Config.ClientSecret)
				assert.Nil(t, ma.AccessToken)
				assert.Nil(t, ma.TokenIV)
			},
		},
		{
			name:         "creates mal auth with long client credentials",
			clientID:     "very-long-client-id-that-might-be-used-in-production",
			clientSecret: "very-long-client-secret-that-should-be-kept-secure",
			accessToken:  []byte("very-long-access-token-string"),
			tokenIV:      []byte("very-long-token-iv-string"),
			validate: func(t *testing.T, ma *MalAuth) {
				assert.Equal(t, "very-long-client-id-that-might-be-used-in-production", ma.Config.ClientID)
				assert.Equal(t, "very-long-client-secret-that-should-be-kept-secure", ma.Config.ClientSecret)
				assert.Equal(t, []byte("very-long-access-token-string"), ma.AccessToken)
				assert.Equal(t, []byte("very-long-token-iv-string"), ma.TokenIV)
			},
		},
		{
			name:         "creates mal auth with special characters in credentials",
			clientID:     "client-id-with-special-chars-!@#$%",
			clientSecret: "secret-with-special-chars-!@#$%",
			accessToken:  []byte("token-with-special-!@#$%"),
			tokenIV:      []byte("iv-with-special-!@#$%"),
			validate: func(t *testing.T, ma *MalAuth) {
				assert.Equal(t, "client-id-with-special-chars-!@#$%", ma.Config.ClientID)
				assert.Equal(t, "secret-with-special-chars-!@#$%", ma.Config.ClientSecret)
				assert.Equal(t, []byte("token-with-special-!@#$%"), ma.AccessToken)
				assert.Equal(t, []byte("iv-with-special-!@#$%"), ma.TokenIV)
			},
		},
		{
			name:         "verifies OAuth endpoint URLs are correct",
			clientID:     "test-id",
			clientSecret: "test-secret",
			accessToken:  []byte("token"),
			tokenIV:      []byte("iv"),
			validate: func(t *testing.T, ma *MalAuth) {
				assert.Equal(t, "https://myanimelist.net/v1/oauth2/authorize", ma.Config.Endpoint.AuthURL)
				assert.Equal(t, "https://myanimelist.net/v1/oauth2/token", ma.Config.Endpoint.TokenURL)
				assert.Equal(t, oauth2.AuthStyleInParams, ma.Config.Endpoint.AuthStyle)
			},
		},
		{
			name:         "verifies ID is always set to 1",
			clientID:     "id1",
			clientSecret: "secret1",
			accessToken:  []byte("token1"),
			tokenIV:      []byte("iv1"),
			validate: func(t *testing.T, ma *MalAuth) {
				assert.Equal(t, 1, ma.Id)
			},
		},
		{
			name:         "creates mal auth with binary token data",
			clientID:     "client",
			clientSecret: "secret",
			accessToken:  []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
			tokenIV:      []byte{0xAA, 0xBB, 0xCC, 0xDD},
			validate: func(t *testing.T, ma *MalAuth) {
				assert.Equal(t, []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}, ma.AccessToken)
				assert.Equal(t, []byte{0xAA, 0xBB, 0xCC, 0xDD}, ma.TokenIV)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewMalAuth(tt.clientID, tt.clientSecret, tt.accessToken, tt.tokenIV)
			require.NotNil(t, result)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// Test that constants are correctly defined
func TestMalAuthConstants(t *testing.T) {
	assert.Equal(t, "https://myanimelist.net/v1/oauth2/authorize", string(AuthURL))
	assert.Equal(t, "https://myanimelist.net/v1/oauth2/token", string(TokenURL))
}

