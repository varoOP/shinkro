package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/varoOP/shinkro/pkg/plex"
)

func TestHamaMalAgent(t *testing.T) {
	tests := []struct {
		name  string
		guid  string
		agent string
		want1 string
		want2 int
		want3 error
	}{
		{
			name:  "HamaTVTVDB",
			guid:  "com.plexapp.agents.hama://tvdb-81797/4/1?lang=en",
			agent: "hama",
			want1: "tvdb",
			want2: 81797,
			want3: nil,
		},
		{
			name:  "HamaTVAniDB",
			guid:  "com.plexapp.agents.hama://anidb-17449/1/4?lang=en",
			agent: "hama",
			want1: "anidb",
			want2: 17449,
			want3: nil,
		},
		{
			name:  "HAMAMovieAniDB",
			guid:  "com.plexapp.agents.hama://anidb-5693?lang=en",
			agent: "hama",
			want1: "anidb",
			want2: 5693,
			want3: nil,
		},
		{
			name:  "HAMAMovieTMDB",
			guid:  "com.plexapp.agents.hama://tmdb-23150?lang=en",
			agent: "hama",
			want1: "tmdb",
			want2: 23150,
			want3: nil,
		},
		{
			name:  "MALTV",
			guid:  "net.fribbtastic.coding.plex.myanimelist://52305/1/1?lang=en",
			agent: "mal",
			want1: "myanimelist",
			want2: 52305,
			want3: nil,
		},
		{
			name:  "MALMovie",
			guid:  "net.fribbtastic.coding.plex.myanimelist://2593?lang=en",
			agent: "mal",
			want1: "myanimelist",
			want2: 2593,
			want3: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, got3 := hamaMALAgent(tt.guid, tt.agent)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func TestPlexAgent(t *testing.T) {
	tests := []struct {
		name      string
		guid      plex.GUID
		mediaType string
		want1     string
		want2     int
		want3     error
	}{
		{
			name: "PlexTV",
			guid: plex.GUID{GUIDS: []struct {
				ID string "json:\"id\""
			}{{ID: "imdb://tt21210326"}, {ID: "tmdb://205308"}, {ID: "tvdb://421994"}},
				GUID: "plex://show/63031ea849f1f16d698849ba"},
			mediaType: "episode",
			want1:     "tvdb",
			want2:     421994,
			want3:     nil,
		},
		{
			name: "PlexMovie",
			guid: plex.GUID{GUIDS: []struct {
				ID string "json:\"id\""
			}{{ID: "imdb://tt0259534"}, {ID: "tmdb://84092"}, {ID: "tvdb://64694"}},
				GUID: "plex://movie/5d7768d196b655001fdc2678"},
			mediaType: "movie",
			want1:     "tmdb",
			want2:     84092,
			want3:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, got3 := plexAgent(tt.guid, tt.mediaType)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}
