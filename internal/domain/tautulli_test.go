package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTautulli(t *testing.T) {
	tests := []struct {
		name        string
		payload     []byte
		expectedErr bool
		validate    func(*testing.T, *Tautulli)
	}{
		{
			name: "valid tautulli payload for episode",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
					"index": "5",
					"librarySectionTitle": "Anime",
					"parentIndex": "2",
					"title": "Episode 5",
					"type": "episode"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, tautulli *Tautulli) {
				assert.Equal(t, PlexScrobbleEvent, tautulli.Event)
				assert.Equal(t, "TestUser", tautulli.Account.Title)
				assert.Equal(t, "Test Anime", tautulli.Metadata.GrandparentTitle)
				assert.Equal(t, "5", tautulli.Metadata.Index)
				assert.Equal(t, "2", tautulli.Metadata.ParentIndex)
				assert.Equal(t, PlexEpisode, tautulli.Metadata.Type)
			},
		},
		{
			name: "valid tautulli payload for movie",
			payload: []byte(`{
				"event": "media.rate",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "",
					"grandparentTitle": "",
					"guid": "net.fribbtastic.coding.plex.myanimelist://12345?lang=en",
					"index": "1",
					"librarySectionTitle": "Anime Movies",
					"parentIndex": "1",
					"title": "Test Movie",
					"type": "movie"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, tautulli *Tautulli) {
				assert.Equal(t, PlexRateEvent, tautulli.Event)
				assert.Equal(t, PlexMovie, tautulli.Metadata.Type)
				assert.Equal(t, "Test Movie", tautulli.Metadata.Title)
			},
		},
		{
			name: "valid tautulli payload with GUID array",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": [
						{"id": "tvdb://12345"},
						{"id": "tmdb://67890"}
					],
					"index": "10",
					"librarySectionTitle": "Anime",
					"parentIndex": "1",
					"title": "Episode 10",
					"type": "episode"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, tautulli *Tautulli) {
				assert.Equal(t, "10", tautulli.Metadata.Index)
				assert.NotNil(t, tautulli.Metadata.GUID.GUIDS)
				assert.Equal(t, 2, len(tautulli.Metadata.GUID.GUIDS))
			},
		},
		{
			name:        "invalid JSON payload",
			payload:     []byte(`{invalid json}`),
			expectedErr: true,
		},
		{
			name:        "empty payload",
			payload:     []byte(`{}`),
			expectedErr: false,
			validate: func(t *testing.T, tautulli *Tautulli) {
				assert.NotNil(t, tautulli)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTautulli(tt.payload)
			if tt.expectedErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestToPlex(t *testing.T) {
	tests := []struct {
		name        string
		payload     []byte
		expectedErr bool
		errContains string
		validate    func(*testing.T, *Plex)
	}{
		{
			name: "converts tautulli episode to plex",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
					"index": "5",
					"librarySectionTitle": "Anime",
					"parentIndex": "2",
					"title": "Episode 5",
					"type": "episode"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, plex *Plex) {
				assert.Equal(t, TautulliWebhook, plex.Source)
				assert.Equal(t, PlexScrobbleEvent, plex.Event)
				assert.Equal(t, "TestUser", plex.Account.Title)
				assert.Equal(t, "Test Anime", plex.Metadata.GrandparentTitle)
				assert.Equal(t, 5, plex.Metadata.Index)
				assert.Equal(t, 2, plex.Metadata.ParentIndex)
				assert.Equal(t, PlexEpisode, plex.Metadata.Type)
				assert.False(t, plex.TimeStamp.IsZero())
			},
		},
		{
			name: "converts tautulli movie to plex",
			payload: []byte(`{
				"event": "media.rate",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "",
					"grandparentTitle": "",
					"guid": "net.fribbtastic.coding.plex.myanimelist://12345?lang=en",
					"index": "1",
					"librarySectionTitle": "Anime Movies",
					"parentIndex": "1",
					"title": "Test Movie",
					"type": "movie"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, plex *Plex) {
				assert.Equal(t, TautulliWebhook, plex.Source)
				assert.Equal(t, PlexRateEvent, plex.Event)
				assert.Equal(t, PlexMovie, plex.Metadata.Type)
				assert.Equal(t, 1, plex.Metadata.Index)
				assert.Equal(t, 1, plex.Metadata.ParentIndex)
			},
		},
		{
			name: "converts with zero index values",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
					"index": "0",
					"librarySectionTitle": "Anime",
					"parentIndex": "0",
					"title": "Episode 0",
					"type": "episode"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, plex *Plex) {
				assert.Equal(t, 0, plex.Metadata.Index)
				assert.Equal(t, 0, plex.Metadata.ParentIndex)
			},
		},
		{
			name: "converts with large index values",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
					"index": "999",
					"librarySectionTitle": "Anime",
					"parentIndex": "50",
					"title": "Episode 999",
					"type": "episode"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, plex *Plex) {
				assert.Equal(t, 999, plex.Metadata.Index)
				assert.Equal(t, 50, plex.Metadata.ParentIndex)
			},
		},
		{
			name:        "returns error for invalid JSON",
			payload:     []byte(`{invalid json}`),
			expectedErr: true,
			errContains: "invalid character",
		},
		{
			name: "returns error for non-numeric index",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
					"index": "abc",
					"librarySectionTitle": "Anime",
					"parentIndex": "2",
					"title": "Episode 5",
					"type": "episode"
				}
			}`),
			expectedErr: true,
			errContains: "invalid syntax",
		},
		{
			name: "returns error for non-numeric parentIndex",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
					"index": "5",
					"librarySectionTitle": "Anime",
					"parentIndex": "xyz",
					"title": "Episode 5",
					"type": "episode"
				}
			}`),
			expectedErr: true,
			errContains: "invalid syntax",
		},
		{
			name: "handles empty index strings",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
					"index": "",
					"librarySectionTitle": "Anime",
					"parentIndex": "",
					"title": "Episode 5",
					"type": "episode"
				}
			}`),
			expectedErr: true,
			errContains: "invalid syntax",
		},
		{
			name: "preserves GUID structure",
			payload: []byte(`{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"grandparentKey": "/library/metadata/123",
					"grandparentTitle": "Test Anime",
					"guid": [
						{"id": "tvdb://12345"},
						{"id": "tmdb://67890"}
					],
					"index": "5",
					"librarySectionTitle": "Anime",
					"parentIndex": "2",
					"title": "Episode 5",
					"type": "episode"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, plex *Plex) {
				assert.NotNil(t, plex.Metadata.GUID.GUIDS)
				assert.Equal(t, 2, len(plex.Metadata.GUID.GUIDS))
				assert.Equal(t, "tvdb://12345", plex.Metadata.GUID.GUIDS[0].ID)
				assert.Equal(t, "tmdb://67890", plex.Metadata.GUID.GUIDS[1].ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToPlex(tt.payload)
			if tt.expectedErr {
				require.Error(t, err)
				assert.Nil(t, result)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

// Test that ToPlex properly handles the conversion flow
func TestToPlex_Integration(t *testing.T) {
	// Test that NewTautulli + ToPlex works end-to-end
	payload := []byte(`{
		"event": "media.scrobble",
		"Account": {
			"title": "IntegrationTestUser"
		},
		"Metadata": {
			"grandparentKey": "/library/metadata/456",
			"grandparentTitle": "Integration Test Anime",
					"guid": "com.plexapp.agents.hama://anidb-99999/3/7?lang=en",
			"index": "7",
			"librarySectionTitle": "Anime",
			"parentIndex": "3",
			"title": "Episode 7",
			"type": "episode"
		}
	}`)

	// First create Tautulli
	tautulli, err := NewTautulli(payload)
	require.NoError(t, err)
	require.NotNil(t, tautulli)

	// Then convert to Plex
	plex, err := ToPlex(payload)
	require.NoError(t, err)
	require.NotNil(t, plex)

	// Verify conversion
	assert.Equal(t, TautulliWebhook, plex.Source)
	assert.Equal(t, tautulli.Event, plex.Event)
	assert.Equal(t, tautulli.Account.Title, plex.Account.Title)
	assert.Equal(t, tautulli.Metadata.GrandparentTitle, plex.Metadata.GrandparentTitle)
	assert.Equal(t, 7, plex.Metadata.Index)
	assert.Equal(t, 3, plex.Metadata.ParentIndex)
}

// Test edge case: negative index values (should fail conversion)
func TestToPlex_NegativeIndex(t *testing.T) {
	payload := []byte(`{
		"event": "media.scrobble",
		"Account": {
			"title": "TestUser"
		},
		"Metadata": {
			"grandparentKey": "/library/metadata/123",
			"grandparentTitle": "Test Anime",
			"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
			"index": "-1",
			"librarySectionTitle": "Anime",
			"parentIndex": "2",
			"title": "Episode -1",
			"type": "episode"
		}
	}`)

	// Should succeed in parsing (negative numbers are valid integers)
	plex, err := ToPlex(payload)
	require.NoError(t, err)
	assert.Equal(t, -1, plex.Metadata.Index)
}
