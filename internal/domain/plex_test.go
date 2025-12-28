package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlexWebhook(t *testing.T) {
	tests := []struct {
		name        string
		payload     []byte
		expectedErr bool
		validate    func(*testing.T, *Plex)
	}{
		{
			name: "valid plex webhook payload",
			payload: []byte(`{
				"event": "media.scrobble",
				"rating": 8.0,
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"type": "episode",
					"title": "Episode 1",
					"grandparentTitle": "Test Anime",
					"index": 1,
					"parentIndex": 1,
					"librarySectionTitle": "Anime",
					"guid": "com.plexapp.agents.hama://anidb-12345/1/1?lang=en"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, p *Plex) {
				assert.Equal(t, PlexWebhook, p.Source)
				assert.Equal(t, PlexScrobbleEvent, p.Event)
				assert.Equal(t, float32(8.0), p.Rating)
				assert.Equal(t, "TestUser", p.Account.Title)
				assert.Equal(t, PlexEpisode, p.Metadata.Type)
				assert.False(t, p.TimeStamp.IsZero())
			},
		},
		{
			name: "valid rate event payload",
			payload: []byte(`{
				"event": "media.rate",
				"rating": 9.5,
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"type": "movie",
					"title": "Test Movie",
					"librarySectionTitle": "Anime Movies",
					"guid": "net.fribbtastic.coding.plex.myanimelist://12345?lang=en"
				}
			}`),
			expectedErr: false,
			validate: func(t *testing.T, p *Plex) {
				assert.Equal(t, PlexRateEvent, p.Event)
				assert.Equal(t, float32(9.5), p.Rating)
				assert.Equal(t, PlexMovie, p.Metadata.Type)
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
			validate: func(t *testing.T, p *Plex) {
				assert.Equal(t, PlexWebhook, p.Source)
				assert.False(t, p.TimeStamp.IsZero())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewPlexWebhook(tt.payload)
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

func TestPlex_GetPlexMediaType(t *testing.T) {
	tests := []struct {
		name     string
		plex     *Plex
		expected PlexMediaType
	}{
		{
			name: "returns episode for episode type",
			plex: &Plex{
				Metadata: Metadata{
					Type: PlexEpisode,
				},
			},
			expected: PlexEpisode,
		},
		{
			name: "returns movie for movie type",
			plex: &Plex{
				Metadata: Metadata{
					Type: PlexMovie,
				},
			},
			expected: PlexMovie,
		},
		{
			name: "returns empty string for unknown type",
			plex: &Plex{
				Metadata: Metadata{
					Type: "",
				},
			},
			expected: "",
		},
		{
			name: "returns empty string for invalid type",
			plex: &Plex{
				Metadata: Metadata{
					Type: "invalid",
				},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plex.GetPlexMediaType()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlex_GetPlexEvent(t *testing.T) {
	tests := []struct {
		name     string
		plex     *Plex
		expected PlexEvent
	}{
		{
			name: "returns scrobble event",
			plex: &Plex{
				Event: PlexScrobbleEvent,
			},
			expected: PlexScrobbleEvent,
		},
		{
			name: "returns rate event",
			plex: &Plex{
				Event: PlexRateEvent,
			},
			expected: PlexRateEvent,
		},
		{
			name: "returns empty string for unknown event",
			plex: &Plex{
				Event: "",
			},
			expected: "",
		},
		{
			name: "returns empty string for invalid event",
			plex: &Plex{
				Event: "invalid.event",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plex.GetPlexEvent()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlex_IsEventAllowed(t *testing.T) {
	tests := []struct {
		name     string
		plex     *Plex
		expected bool
	}{
		{
			name: "allows scrobble event",
			plex: &Plex{
				Event: PlexScrobbleEvent,
			},
			expected: true,
		},
		{
			name: "allows rate event",
			plex: &Plex{
				Event: PlexRateEvent,
			},
			expected: true,
		},
		{
			name: "disallows unknown event",
			plex: &Plex{
				Event: "",
			},
			expected: false,
		},
		{
			name: "disallows invalid event",
			plex: &Plex{
				Event: "media.play",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plex.IsEventAllowed()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlex_IsRatingAllowed(t *testing.T) {
	tests := []struct {
		name     string
		plex     *Plex
		expected bool
	}{
		{
			name: "allows positive rating",
			plex: &Plex{
				Rating: 8.5,
			},
			expected: true,
		},
		{
			name: "allows zero rating",
			plex: &Plex{
				Rating: 0,
			},
			expected: true,
		},
		{
			name: "allows maximum rating",
			plex: &Plex{
				Rating: 10.0,
			},
			expected: true,
		},
		{
			name: "disallows negative rating",
			plex: &Plex{
				Rating: -1.0,
			},
			expected: false,
		},
		{
			name: "disallows very negative rating",
			plex: &Plex{
				Rating: -10.0,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plex.IsRatingAllowed()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlex_IsPlexUserAllowed(t *testing.T) {
	tests := []struct {
		name     string
		plex     *Plex
		settings *PlexSettings
		expected bool
	}{
		{
			name: "allows matching user",
			plex: &Plex{
				Account: struct {
					Id           int    `json:"id"`
					ThumbnailUrl string `json:"thumb"`
					Title        string `json:"title"`
				}{
					Title: "TestUser",
				},
			},
			settings: &PlexSettings{
				PlexUser: "TestUser",
			},
			expected: true,
		},
		{
			name: "disallows non-matching user",
			plex: &Plex{
				Account: struct {
					Id           int    `json:"id"`
					ThumbnailUrl string `json:"thumb"`
					Title        string `json:"title"`
				}{
					Title: "TestUser",
				},
			},
			settings: &PlexSettings{
				PlexUser: "OtherUser",
			},
			expected: false,
		},
		{
			name: "disallows empty user when settings has user",
			plex: &Plex{
				Account: struct {
					Id           int    `json:"id"`
					ThumbnailUrl string `json:"thumb"`
					Title        string `json:"title"`
				}{
					Title: "",
				},
			},
			settings: &PlexSettings{
				PlexUser: "TestUser",
			},
			expected: false,
		},
		{
			name: "allows user when both are empty",
			plex: &Plex{
				Account: struct {
					Id           int    `json:"id"`
					ThumbnailUrl string `json:"thumb"`
					Title        string `json:"title"`
				}{
					Title: "",
				},
			},
			settings: &PlexSettings{
				PlexUser: "",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plex.IsPlexUserAllowed(tt.settings)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlex_IsAnimeLibrary(t *testing.T) {
	tests := []struct {
		name     string
		plex     *Plex
		settings *PlexSettings
		expected bool
	}{
		{
			name: "allows matching library",
			plex: &Plex{
				Metadata: Metadata{
					LibrarySectionTitle: "Anime",
				},
			},
			settings: &PlexSettings{
				AnimeLibraries: []string{"Anime"},
			},
			expected: true,
		},
		{
			name: "allows matching library from multiple",
			plex: &Plex{
				Metadata: Metadata{
					LibrarySectionTitle: "Anime Movies",
				},
			},
			settings: &PlexSettings{
				AnimeLibraries: []string{"Anime", "Anime Movies", "Other"},
			},
			expected: true,
		},
		{
			name: "disallows non-matching library",
			plex: &Plex{
				Metadata: Metadata{
					LibrarySectionTitle: "TV Shows",
				},
			},
			settings: &PlexSettings{
				AnimeLibraries: []string{"Anime", "Anime Movies"},
			},
			expected: false,
		},
		{
			name: "disallows when libraries list is empty",
			plex: &Plex{
				Metadata: Metadata{
					LibrarySectionTitle: "Anime",
				},
			},
			settings: &PlexSettings{
				AnimeLibraries: []string{},
			},
			expected: false,
		},
		{
			name: "disallows when library section title is empty",
			plex: &Plex{
				Metadata: Metadata{
					LibrarySectionTitle: "",
				},
			},
			settings: &PlexSettings{
				AnimeLibraries: []string{"Anime"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plex.IsAnimeLibrary(tt.settings)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlex_IsMediaTypeAllowed(t *testing.T) {
	tests := []struct {
		name     string
		plex     *Plex
		expected bool
	}{
		{
			name: "allows episode type",
			plex: &Plex{
				Metadata: Metadata{
					Type: PlexEpisode,
				},
			},
			expected: true,
		},
		{
			name: "allows movie type",
			plex: &Plex{
				Metadata: Metadata{
					Type: PlexMovie,
				},
			},
			expected: true,
		},
		{
			name: "disallows empty type",
			plex: &Plex{
				Metadata: Metadata{
					Type: "",
				},
			},
			expected: false,
		},
		{
			name: "disallows invalid type",
			plex: &Plex{
				Metadata: Metadata{
					Type: "show",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plex.IsMediaTypeAllowed()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlex_SetAnimeFields(t *testing.T) {
	tests := []struct {
		name     string
		plex     *Plex
		source   PlexSupportedDBs
		id       int
		validate func(*testing.T, AnimeUpdate)
	}{
		{
			name: "sets fields for movie",
			plex: &Plex{
				ID: 12345,
				Metadata: Metadata{
					Type:  PlexMovie,
					Title: "Test Movie",
				},
			},
			source: TMDB,
			id:     67890,
			validate: func(t *testing.T, au AnimeUpdate) {
				assert.Equal(t, int64(12345), au.PlexId)
				assert.Equal(t, TMDB, au.SourceDB)
				assert.Equal(t, 67890, au.SourceId)
				assert.Equal(t, 1, au.SeasonNum)
				assert.Equal(t, 1, au.EpisodeNum)
				assert.Equal(t, "Test Movie", au.ListDetails.Title)
				assert.NotZero(t, au.Timestamp)
				assert.NotNil(t, au.Plex)
			},
		},
		{
			name: "sets fields for episode",
			plex: &Plex{
				ID: 12345,
				Metadata: Metadata{
					Type:             PlexEpisode,
					GrandparentTitle: "Test Anime",
					Index:            5,
					ParentIndex:      2,
				},
			},
			source: TVDB,
			id:     11111,
			validate: func(t *testing.T, au AnimeUpdate) {
				assert.Equal(t, int64(12345), au.PlexId)
				assert.Equal(t, TVDB, au.SourceDB)
				assert.Equal(t, 11111, au.SourceId)
				assert.Equal(t, 2, au.SeasonNum)
				assert.Equal(t, 5, au.EpisodeNum)
				assert.Equal(t, "Test Anime", au.ListDetails.Title)
				assert.NotZero(t, au.Timestamp)
				assert.NotNil(t, au.Plex)
			},
		},
		{
			name: "sets fields for episode with AniDB source",
			plex: &Plex{
				ID: 99999,
				Metadata: Metadata{
					Type:             PlexEpisode,
					GrandparentTitle: "Another Anime",
					Index:            12,
					ParentIndex:      1,
				},
			},
			source: AniDB,
			id:     22222,
			validate: func(t *testing.T, au AnimeUpdate) {
				assert.Equal(t, AniDB, au.SourceDB)
				assert.Equal(t, 22222, au.SourceId)
				assert.Equal(t, 1, au.SeasonNum)
				assert.Equal(t, 12, au.EpisodeNum)
				assert.Equal(t, "Another Anime", au.ListDetails.Title)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plex.SetAnimeFields(tt.source, tt.id)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestPlex_IsMetadataAgentAllowed(t *testing.T) {
	tests := []struct {
		name          string
		plex          *Plex
		expectedBool  bool
		expectedAgent PlexSupportedAgents
	}{
		{
			name: "detects HAMA agent from GUID",
			plex: &Plex{
				Metadata: Metadata{
					GUID: GUID{
						GUID: "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
					},
				},
			},
			expectedBool:  true,
			expectedAgent: HAMA,
		},
		{
			name: "detects HAMA agent from GrandparentGUID",
			plex: &Plex{
				Metadata: Metadata{
					GrandparentGUID: "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
				},
			},
			expectedBool:  true,
			expectedAgent: HAMA,
		},
		{
			name: "detects MAL agent from GUID",
			plex: &Plex{
				Metadata: Metadata{
					GUID: GUID{
						GUID: "net.fribbtastic.coding.plex.myanimelist://12345?lang=en",
					},
				},
			},
			expectedBool:  true,
			expectedAgent: MALAgent,
		},
		{
			name: "detects Plex agent from GUID",
			plex: &Plex{
				Metadata: Metadata{
					GUID: GUID{
						GUID: "plex://tvdb-12345",
					},
				},
			},
			expectedBool:  true,
			expectedAgent: PlexAgent,
		},
		{
			name: "returns false for unsupported agent",
			plex: &Plex{
				Metadata: Metadata{
					GUID: GUID{
						GUID: "com.plexapp.agents.thetvdb://12345",
					},
				},
			},
			expectedBool:  false,
			expectedAgent: "",
		},
		{
			name: "returns false for empty GUID",
			plex: &Plex{
				Metadata: Metadata{
					GUID: GUID{
						GUID: "",
					},
				},
			},
			expectedBool:  false,
			expectedAgent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultBool, resultAgent := tt.plex.IsMetadataAgentAllowed()
			assert.Equal(t, tt.expectedBool, resultBool)
			assert.Equal(t, tt.expectedAgent, resultAgent)
		})
	}
}

func TestGUID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectedErr bool
		validate    func(*testing.T, *GUID)
	}{
		{
			name:        "unmarshals string GUID",
			data:        []byte(`"com.plexapp.agents.hama://anidb-12345/1/1"`),
			expectedErr: false,
			validate: func(t *testing.T, g *GUID) {
				assert.Equal(t, "com.plexapp.agents.hama://anidb-12345/1/1", g.GUID)
				assert.Nil(t, g.GUIDS)
			},
		},
		{
			name:        "unmarshals array of GUIDs",
			data:        []byte(`[{"id":"tvdb://12345"},{"id":"tmdb://67890"}]`),
			expectedErr: false,
			validate: func(t *testing.T, g *GUID) {
				assert.NotNil(t, g.GUIDS)
				assert.Equal(t, 2, len(g.GUIDS))
				assert.Equal(t, "tvdb://12345", g.GUIDS[0].ID)
				assert.Equal(t, "tmdb://67890", g.GUIDS[1].ID)
			},
		},
		{
			name:        "unmarshals empty string",
			data:        []byte(`""`),
			expectedErr: false,
			validate: func(t *testing.T, g *GUID) {
				assert.Equal(t, "", g.GUID)
			},
		},
		{
			name:        "unmarshals empty array",
			data:        []byte(`[]`),
			expectedErr: false,
			validate: func(t *testing.T, g *GUID) {
				assert.NotNil(t, g.GUIDS)
				assert.Equal(t, 0, len(g.GUIDS))
			},
		},
		{
			name:        "returns error for invalid JSON",
			data:        []byte(`{invalid}`),
			expectedErr: true,
		},
		{
			name:        "returns error for number",
			data:        []byte(`12345`),
			expectedErr: true,
		},
		{
			name:        "returns error for boolean",
			data:        []byte(`true`),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GUID{}
			err := json.Unmarshal(tt.data, g)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, g)
				}
			}
		})
	}
}

func TestGUID_HamaMALAgent(t *testing.T) {
	tests := []struct {
		name        string
		guid        *GUID
		agent       PlexSupportedAgents
		expectedDB  PlexSupportedDBs
		expectedID  int
		expectedErr bool
		errContains string
	}{
		{
			name: "parses HAMA agent with anidb",
			guid: &GUID{
				GUID: "com.plexapp.agents.hama://anidb-12345/1/1?lang=en",
			},
			agent:       HAMA,
			expectedDB:  AniDB,
			expectedID:  12345,
			expectedErr: false,
		},
		{
			name: "parses HAMA agent with tvdb",
			guid: &GUID{
				GUID: "com.plexapp.agents.hama://tvdb-67890/2/5?lang=en",
			},
			agent:       HAMA,
			expectedDB:  TVDB,
			expectedID:  67890,
			expectedErr: false,
		},
		{
			name: "parses MAL agent",
			guid: &GUID{
				GUID: "net.fribbtastic.coding.plex.myanimelist://12345?lang=en",
			},
			agent:       MALAgent,
			expectedDB:  MAL,
			expectedID:  12345,
			expectedErr: false,
		},
		{
			name: "returns error for invalid GUID format",
			guid: &GUID{
				GUID: "invalid-guid-format",
			},
			agent:       HAMA,
			expectedErr: true,
			errContains: "unable to parse GUID",
		},
		{
			name: "returns error for non-numeric ID - regex doesn't match",
			guid: &GUID{
				GUID: "com.plexapp.agents.hama://anidb-abc/1/1",
			},
			agent:       HAMA,
			expectedErr: true,
			errContains: "unable to parse GUID",
		},
		{
			name: "returns error for unsupported agent type",
			guid: &GUID{
				GUID: "com.plexapp.agents.hama://anidb-12345/1/1",
			},
			agent:       PlexAgent,
			expectedErr: true,
			errContains: "unsupported agent type for HamaMALAgent",
		},
		{
			name: "handles GUID with spaces",
			guid: &GUID{
				GUID: "com.plexapp.agents.hama://anidb -12345/1/1",
			},
			agent:       HAMA,
			expectedDB:  AniDB,
			expectedID:  12345,
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultDB, resultID, err := tt.guid.HamaMALAgent(tt.agent)
			if tt.expectedErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, PlexSupportedDBs(""), resultDB)
				assert.Equal(t, -1, resultID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedDB, resultDB)
				assert.Equal(t, tt.expectedID, resultID)
			}
		})
	}
}

func TestGUID_PlexAgent(t *testing.T) {
	tests := []struct {
		name        string
		guid        *GUID
		mediaType   PlexMediaType
		expectedDB  PlexSupportedDBs
		expectedID  int
		expectedErr bool
		errContains string
	}{
		{
			name: "parses TVDB GUID for episode",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{
					{ID: "tvdb://12345"},
					{ID: "tmdb://67890"},
				},
			},
			mediaType:   PlexEpisode,
			expectedDB:  TVDB,
			expectedID:  12345,
			expectedErr: false,
		},
		{
			name: "parses TMDB GUID for movie",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{
					{ID: "tvdb://12345"},
					{ID: "tmdb://67890"},
				},
			},
			mediaType:   PlexMovie,
			expectedDB:  TMDB,
			expectedID:  67890,
			expectedErr: false,
		},
		{
			name: "returns first matching TVDB for episode",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{
					{ID: "tvdb://11111"},
					{ID: "tvdb://22222"},
					{ID: "tmdb://33333"},
				},
			},
			mediaType:   PlexEpisode,
			expectedDB:  TVDB,
			expectedID:  11111,
			expectedErr: false,
		},
		{
			name: "returns error for non-numeric ID",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{
					{ID: "tvdb://abc"},
				},
			},
			mediaType:   PlexEpisode,
			expectedErr: true,
			errContains: "id conversion failed",
		},
		{
			name: "returns error when no matching database found for episode",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{
					{ID: "tmdb://12345"},
				},
			},
			mediaType:   PlexEpisode,
			expectedErr: true,
			errContains: "no supported online database found",
		},
		{
			name: "returns error when no matching database found for movie",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{
					{ID: "tvdb://12345"},
				},
			},
			mediaType:   PlexMovie,
			expectedErr: true,
			errContains: "no supported online database found",
		},
		{
			name: "returns error for empty GUIDS array",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{},
			},
			mediaType:   PlexEpisode,
			expectedErr: true,
			errContains: "no supported online database found",
		},
		{
			name: "returns error for invalid GUID format",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{
					{ID: "invalid-format"},
				},
			},
			mediaType:   PlexEpisode,
			expectedErr: true,
			errContains: "no supported online database found",
		},
		{
			name: "handles GUID with multiple colons",
			guid: &GUID{
				GUIDS: []struct {
					ID string `json:"id"`
				}{
					{ID: "tvdb://12345://extra"},
				},
			},
			mediaType:   PlexEpisode,
			expectedDB:  TVDB,
			expectedID:  12345,
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultDB, resultID, err := tt.guid.PlexAgent(tt.mediaType)
			if tt.expectedErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, PlexSupportedDBs(""), resultDB)
				assert.Equal(t, -1, resultID)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedDB, resultDB)
				assert.Equal(t, tt.expectedID, resultID)
			}
		})
	}
}
