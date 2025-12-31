package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTautulli(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		expectedError bool
		validate      func(*testing.T, *Tautulli)
	}{
		{
			name:          "valid episode payload",
			payload:       `{"Account":{"title":"TestUser"},"event":"media.scrobble","Metadata":{"title":"Episode 5","type":"episode","parentIndex":"1","index":"5","guid":"com.plexapp.agents.hama://anidb-12345/1/1?lang=en","grandparentKey":"/library/metadata/123","grandparentTitle":"Attack on Titan","librarySectionTitle":"Anime"}}`,
			expectedError: false,
			validate: func(t *testing.T, tautulli *Tautulli) {
				assert.Equal(t, "TestUser", tautulli.Account.Title)
				assert.Equal(t, PlexScrobbleEvent, tautulli.Event)
				assert.Equal(t, PlexEpisode, tautulli.Metadata.Type)
				assert.Equal(t, "Attack on Titan", tautulli.Metadata.GrandparentTitle)
				assert.Equal(t, "5", tautulli.Metadata.Index)
				assert.Equal(t, "1", tautulli.Metadata.ParentIndex)
			},
		},
		{
			name:          "valid movie payload",
			payload:       `{"Account":{"title":"TestUser"},"event":"media.rate","Metadata":{"title":"Your Name","type":"movie","parentIndex":"1","index":"1","guid":[{"id":"tmdb://372058"}],"grandparentKey":"","grandparentTitle":"","librarySectionTitle":"Anime Movies"}}`,
			expectedError: false,
			validate: func(t *testing.T, tautulli *Tautulli) {
				assert.Equal(t, PlexMovie, tautulli.Metadata.Type)
				assert.Equal(t, "Your Name", tautulli.Metadata.Title)
			},
		},
		{
			name:          "invalid JSON",
			payload:       `{"invalid": json}`,
			expectedError: true,
		},
		{
			name:          "empty payload",
			payload:       `{}`,
			expectedError: false,
			validate: func(t *testing.T, tautulli *Tautulli) {
				assert.NotNil(t, tautulli)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewTautulli([]byte(tt.payload))
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestToPlex(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		expectedError bool
		validate      func(*testing.T, *Plex)
	}{
		{
			name:          "valid episode conversion",
			payload:       `{"Account":{"title":"TestUser"},"event":"media.scrobble","Metadata":{"title":"Episode 5","type":"episode","parentIndex":"1","index":"5","guid":"com.plexapp.agents.hama://anidb-12345/1/1?lang=en","grandparentKey":"/library/metadata/123","grandparentTitle":"Attack on Titan","librarySectionTitle":"Anime"}}`,
			expectedError: false,
			validate: func(t *testing.T, plex *Plex) {
				assert.Equal(t, TautulliWebhook, plex.Source)
				assert.Equal(t, PlexScrobbleEvent, plex.Event)
				assert.Equal(t, "TestUser", plex.Account.Title)
				assert.Equal(t, PlexEpisode, plex.Metadata.Type)
				assert.Equal(t, 5, plex.Metadata.Index)
				assert.Equal(t, 1, plex.Metadata.ParentIndex)
				assert.Equal(t, "Attack on Titan", plex.Metadata.GrandparentTitle)
			},
		},
		{
			name:          "valid movie conversion",
			payload:       `{"Account":{"title":"TestUser"},"event":"media.rate","Metadata":{"title":"Your Name","type":"movie","parentIndex":"1","index":"1","guid":[{"id":"tmdb://372058"}],"grandparentKey":"","grandparentTitle":"","librarySectionTitle":"Anime Movies"}}`,
			expectedError: false,
			validate: func(t *testing.T, plex *Plex) {
				assert.Equal(t, TautulliWebhook, plex.Source)
				assert.Equal(t, PlexMovie, plex.Metadata.Type)
				assert.Equal(t, "Your Name", plex.Metadata.Title)
			},
		},
		{
			name:          "invalid JSON",
			payload:       `{"invalid": json}`,
			expectedError: true,
		},
		{
			name:          "invalid parent index",
			payload:       `{"Account":{"title":"Test"},"Metadata":{"parentIndex":"invalid","index":"5","type":"episode"},"event":"media.scrobble"}`,
			expectedError: true,
		},
		{
			name:          "invalid episode index",
			payload:       `{"Account":{"title":"Test"},"Metadata":{"parentIndex":"1","index":"invalid","type":"episode"},"event":"media.scrobble"}`,
			expectedError: true,
		},
		{
			name:          "plex agent with GUID array",
			payload:       `{"Account":{"title":"TestUser"},"event":"media.scrobble","Metadata":{"title":"Episode 12","type":"episode","parentIndex":"2","index":"12","guid":[{"id":"tvdb://362753"},{"id":"tmdb://12345"}],"grandparentKey":"/library/metadata/456","grandparentTitle":"Demon Slayer","librarySectionTitle":"Anime"}}`,
			expectedError: false,
			validate: func(t *testing.T, plex *Plex) {
				assert.Equal(t, TautulliWebhook, plex.Source)
				assert.Equal(t, 12, plex.Metadata.Index)
				assert.Equal(t, 2, plex.Metadata.ParentIndex)
				assert.NotNil(t, plex.Metadata.GUID.GUIDS)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToPlex([]byte(tt.payload))
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestToPlex_IndexConversion(t *testing.T) {
	tests := []struct {
		name          string
		parentIndex   string
		index         string
		expectedError bool
		expectedParent int
		expectedIndex int
	}{
		{
			name:           "valid numeric indices",
			parentIndex:    "5",
			index:          "10",
			expectedError:  false,
			expectedParent: 5,
			expectedIndex:  10,
		},
		{
			name:           "zero indices",
			parentIndex:    "0",
			index:          "0",
			expectedError:  false,
			expectedParent: 0,
			expectedIndex:  0,
		},
		{
			name:           "large indices",
			parentIndex:    "999",
			index:          "1000",
			expectedError:  false,
			expectedParent: 999,
			expectedIndex:  1000,
		},
		{
			name:          "negative parent index",
			parentIndex:   "-1",
			index:         "1",
			expectedError: false,
			expectedParent: -1,
			expectedIndex:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := `{"Account":{"title":"Test"},"Metadata":{"parentIndex":"` + tt.parentIndex + `","index":"` + tt.index + `","type":"episode","grandparentTitle":"Test"},"event":"media.scrobble"}`
			result, err := ToPlex([]byte(payload))
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedParent, result.Metadata.ParentIndex)
				assert.Equal(t, tt.expectedIndex, result.Metadata.Index)
			}
		})
	}
}
