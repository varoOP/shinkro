package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMappingCheckMap(t *testing.T) {
	tests := []struct {
		name       string
		animeMap   *AnimeTVShows
		tvdbid     int
		tvdbseason int
		ep         int
		want1      bool
		want2      int //malid
		want3      int //start
	}{
		{
			name: "DanMachi S4-1",
			animeMap: &AnimeTVShows{
				Anime: []AnimeTV{
					{
						Malid:      47164,
						Tvdbid:     289882,
						TvdbSeason: 4,
						Start:      0,
						UseMapping: false,
					},
					{
						Malid:      53111,
						Tvdbid:     289882,
						TvdbSeason: 4,
						Start:      12,
						UseMapping: false,
					},
				},
			},
			tvdbid:     289882,
			tvdbseason: 4,
			ep:         22,
			want1:      true,
			want2:      53111,
			want3:      12,
		},
		{
			name: "DanMachi S4-2",
			animeMap: &AnimeTVShows{
				Anime: []AnimeTV{
					{
						Malid:      47164,
						Tvdbid:     289882,
						TvdbSeason: 4,
						Start:      0,
						UseMapping: false,
					},
					{
						Malid:      53111,
						Tvdbid:     289882,
						TvdbSeason: 4,
						Start:      12,
						UseMapping: false,
					},
				},
			},
			tvdbid:     289882,
			tvdbseason: 4,
			ep:         10,
			want1:      true,
			want2:      47164,
			want3:      0,
		},
		{
			name: "SPYXFAMILY S1-2",
			animeMap: &AnimeTVShows{
				Anime: []AnimeTV{
					{
						Malid:      50602,
						Tvdbid:     405920,
						TvdbSeason: 1,
						Start:      13,
						UseMapping: false,
					},
					{
						Malid:      50265,
						Tvdbid:     405920,
						TvdbSeason: 1,
						Start:      0,
						UseMapping: false,
					},
				},
			},
			tvdbid:     405920,
			tvdbseason: 1,
			ep:         19,
			want1:      true,
			want2:      50602,
			want3:      13,
		},
		{
			name: "One Piece-1",
			animeMap: &AnimeTVShows{
				Anime: []AnimeTV{
					{
						Malid:      21,
						Tvdbid:     81797,
						TvdbSeason: 0,
						Start:      0,
						UseMapping: true,
						AnimeMapping: []AnimeMapping{
							{
								TvdbSeason: 10,
								Start:      196,
							},
							{
								TvdbSeason: 12,
								Start:      326,
							},
							{
								TvdbSeason: 15,
								Start:      517,
							},
							{
								TvdbSeason: 21,
								Start:      892,
							},
						},
					},
				},
			},
			tvdbid:     81797,
			tvdbseason: 21,
			ep:         186,
			want1:      true,
			want2:      21,
			want3:      892,
		},
		{
			name: "One Piece-2",
			animeMap: &AnimeTVShows{
				Anime: []AnimeTV{
					{
						Malid:      21,
						Tvdbid:     81797,
						TvdbSeason: 0,
						Start:      0,
						UseMapping: true,
						AnimeMapping: []AnimeMapping{
							{
								TvdbSeason: 10,
								Start:      196,
							},
							{
								TvdbSeason: 12,
								Start:      326,
							},
							{
								TvdbSeason: 15,
								Start:      517,
							},
							{
								TvdbSeason: 21,
								Start:      892,
							},
						},
					},
				},
			},
			tvdbid:     81797,
			tvdbseason: 10,
			ep:         30,
			want1:      true,
			want2:      21,
			want3:      196,
		},
		{
			name: "Monogatari S3-1",
			animeMap: &AnimeTVShows{
				Anime: []AnimeTV{
					{
						Malid:      17074,
						Tvdbid:     102261,
						TvdbSeason: 0,
						Start:      0,
						UseMapping: true,
						AnimeMapping: []AnimeMapping{
							{
								TvdbSeason:       0,
								Start:            0,
								MappingType:      "explicit",
								ExplicitEpisodes: map[int]int{7: 6, 8: 11, 9: 16},
							},
							{
								TvdbSeason:      3,
								Start:           1,
								MappingType:     "range",
								SkipMalEpisodes: []int{6, 11, 16},
							},
						},
					},
				},
			},
			tvdbid:     102261,
			tvdbseason: 0,
			ep:         7,
			want1:      true,
			want2:      17074,
			want3:      0,
		},
		{
			name: "Monogatari S3-2",
			animeMap: &AnimeTVShows{
				Anime: []AnimeTV{
					{
						Malid:      17074,
						Tvdbid:     102261,
						TvdbSeason: 0,
						Start:      0,
						UseMapping: true,
						AnimeMapping: []AnimeMapping{
							{
								TvdbSeason:       0,
								Start:            0,
								MappingType:      "explicit",
								ExplicitEpisodes: map[int]int{7: 6, 8: 11, 9: 16},
							},
							{
								TvdbSeason:      3,
								Start:           1,
								MappingType:     "range",
								SkipMalEpisodes: []int{6, 11, 16},
							},
						},
					},
				},
			},
			tvdbid:     102261,
			tvdbseason: 3,
			ep:         23,
			want1:      true,
			want2:      17074,
			want3:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := tt.animeMap.CheckMap(tt.tvdbid, tt.tvdbseason, tt.ep)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2.Malid)
			assert.Equal(t, tt.want3, got2.Start)
		})
	}
}

func TestCalculateEpNumWithExplicitAndSkips(t *testing.T) {
	tests := []struct {
		name        string
		details     *AnimeMapDetails
		tvdbEpisode int
		expectedMal int
	}{
		{
			name: "Monogatari S3-1",
			details: &AnimeMapDetails{
				Malid:            17074,
				Start:            0,
				UseMapping:       true,
				MappingType:      "explicit",
				ExplicitEpisodes: map[int]int{7: 6, 8: 11, 9: 16},
			},
			tvdbEpisode: 7,
			expectedMal: 6,
		},
		{
			name: "Monogatari S3-2",
			details: &AnimeMapDetails{
				Malid:           17074,
				Start:           1,
				UseMapping:      true,
				MappingType:     "range",
				SkipMalEpisodes: []int{6, 11, 16},
			},
			tvdbEpisode: 23,
			expectedMal: 26,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualMal := tt.details.CalculateEpNum(tt.tvdbEpisode)
			assert.Equal(t, tt.expectedMal, actualMal)
		})
	}
}

func TestNewMapSettings(t *testing.T) {
	tests := []struct {
		name     string
		tvdb     bool
		tmdb     bool
		tvdbPath string
		tmdbPath string
		validate func(*testing.T, *MapSettings)
	}{
		{
			name:     "creates map settings with all fields",
			tvdb:     true,
			tmdb:     true,
			tvdbPath: "/path/to/tvdb.yaml",
			tmdbPath: "/path/to/tmdb.yaml",
			validate: func(t *testing.T, ms *MapSettings) {
				assert.True(t, ms.TVDBEnabled)
				assert.True(t, ms.TMDBEnabled)
				assert.Equal(t, "/path/to/tvdb.yaml", ms.CustomMapTVDBPath)
				assert.Equal(t, "/path/to/tmdb.yaml", ms.CustomMapTMDBPath)
			},
		},
		{
			name:     "creates map settings with only TVDB",
			tvdb:     true,
			tmdb:     false,
			tvdbPath: "/path/to/tvdb.yaml",
			tmdbPath: "",
			validate: func(t *testing.T, ms *MapSettings) {
				assert.True(t, ms.TVDBEnabled)
				assert.False(t, ms.TMDBEnabled)
				assert.Equal(t, "/path/to/tvdb.yaml", ms.CustomMapTVDBPath)
				assert.Equal(t, "", ms.CustomMapTMDBPath)
			},
		},
		{
			name:     "creates map settings with only TMDB",
			tvdb:     false,
			tmdb:     true,
			tvdbPath: "",
			tmdbPath: "/path/to/tmdb.yaml",
			validate: func(t *testing.T, ms *MapSettings) {
				assert.False(t, ms.TVDBEnabled)
				assert.True(t, ms.TMDBEnabled)
				assert.Equal(t, "", ms.CustomMapTVDBPath)
				assert.Equal(t, "/path/to/tmdb.yaml", ms.CustomMapTMDBPath)
			},
		},
		{
			name:     "creates map settings with both disabled",
			tvdb:     false,
			tmdb:     false,
			tvdbPath: "",
			tmdbPath: "",
			validate: func(t *testing.T, ms *MapSettings) {
				assert.False(t, ms.TVDBEnabled)
				assert.False(t, ms.TMDBEnabled)
				assert.Equal(t, "", ms.CustomMapTVDBPath)
				assert.Equal(t, "", ms.CustomMapTMDBPath)
			},
		},
		{
			name:     "creates map settings with empty paths",
			tvdb:     true,
			tmdb:     true,
			tvdbPath: "",
			tmdbPath: "",
			validate: func(t *testing.T, ms *MapSettings) {
				assert.True(t, ms.TVDBEnabled)
				assert.True(t, ms.TMDBEnabled)
				assert.Equal(t, "", ms.CustomMapTVDBPath)
				assert.Equal(t, "", ms.CustomMapTMDBPath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewMapSettings(tt.tvdb, tt.tmdb, tt.tvdbPath, tt.tmdbPath)
			require.NotNil(t, result)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestMapSettings_ShouldLoadLocal(t *testing.T) {
	tests := []struct {
		name        string
		settings    *MapSettings
		expectedTVDB bool
		expectedTMDB bool
	}{
		{
			name: "both should load when enabled and paths provided",
			settings: &MapSettings{
				TVDBEnabled:       true,
				TMDBEnabled:       true,
				CustomMapTVDBPath: "/path/to/tvdb.yaml",
				CustomMapTMDBPath: "/path/to/tmdb.yaml",
			},
			expectedTVDB: true,
			expectedTMDB: true,
		},
		{
			name: "TVDB should load, TMDB should not",
			settings: &MapSettings{
				TVDBEnabled:       true,
				TMDBEnabled:       false,
				CustomMapTVDBPath: "/path/to/tvdb.yaml",
				CustomMapTMDBPath: "",
			},
			expectedTVDB: true,
			expectedTMDB: false,
		},
		{
			name: "TMDB should load, TVDB should not",
			settings: &MapSettings{
				TVDBEnabled:       false,
				TMDBEnabled:       true,
				CustomMapTVDBPath: "",
				CustomMapTMDBPath: "/path/to/tmdb.yaml",
			},
			expectedTVDB: false,
			expectedTMDB: true,
		},
		{
			name: "neither should load when both disabled",
			settings: &MapSettings{
				TVDBEnabled:       false,
				TMDBEnabled:       false,
				CustomMapTVDBPath: "",
				CustomMapTMDBPath: "",
			},
			expectedTVDB: false,
			expectedTMDB: false,
		},
		{
			name: "neither should load when enabled but no paths",
			settings: &MapSettings{
				TVDBEnabled:       true,
				TMDBEnabled:       true,
				CustomMapTVDBPath: "",
				CustomMapTMDBPath: "",
			},
			expectedTVDB: false,
			expectedTMDB: false,
		},
		{
			name: "TVDB should not load when enabled but no path",
			settings: &MapSettings{
				TVDBEnabled:       true,
				TMDBEnabled:       false,
				CustomMapTVDBPath: "",
				CustomMapTMDBPath: "",
			},
			expectedTVDB: false,
			expectedTMDB: false,
		},
		{
			name: "TMDB should not load when enabled but no path",
			settings: &MapSettings{
				TVDBEnabled:       false,
				TMDBEnabled:       true,
				CustomMapTVDBPath: "",
				CustomMapTMDBPath: "",
			},
			expectedTVDB: false,
			expectedTMDB: false,
		},
		{
			name: "TVDB should not load when path provided but disabled",
			settings: &MapSettings{
				TVDBEnabled:       false,
				TMDBEnabled:       false,
				CustomMapTVDBPath: "/path/to/tvdb.yaml",
				CustomMapTMDBPath: "",
			},
			expectedTVDB: false,
			expectedTMDB: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tvdb, tmdb := tt.settings.ShouldLoadLocal()
			assert.Equal(t, tt.expectedTVDB, tvdb)
			assert.Equal(t, tt.expectedTMDB, tmdb)
		})
	}
}

func TestAnimeMovies_CheckMap(t *testing.T) {
	tests := []struct {
		name        string
		movies      *AnimeMovies
		tmdbid      int
		expectedFound bool
		expectedMALID int
	}{
		{
			name: "finds movie by TMDB ID",
			movies: &AnimeMovies{
				AnimeMovie: []AnimeMovie{
					{
						MainTitle: "Test Movie 1",
						TMDBID:    12345,
						MALID:     67890,
					},
					{
						MainTitle: "Test Movie 2",
						TMDBID:    11111,
						MALID:     22222,
					},
				},
			},
			tmdbid:        12345,
			expectedFound: true,
			expectedMALID: 67890,
		},
		{
			name: "finds second movie",
			movies: &AnimeMovies{
				AnimeMovie: []AnimeMovie{
					{
						MainTitle: "Test Movie 1",
						TMDBID:    12345,
						MALID:     67890,
					},
					{
						MainTitle: "Test Movie 2",
						TMDBID:    11111,
						MALID:     22222,
					},
				},
			},
			tmdbid:        11111,
			expectedFound: true,
			expectedMALID: 22222,
		},
		{
			name: "returns false when movie not found",
			movies: &AnimeMovies{
				AnimeMovie: []AnimeMovie{
					{
						MainTitle: "Test Movie 1",
						TMDBID:    12345,
						MALID:     67890,
					},
				},
			},
			tmdbid:        99999,
			expectedFound: false,
			expectedMALID: 0,
		},
		{
			name: "returns false for empty movie list",
			movies: &AnimeMovies{
				AnimeMovie: []AnimeMovie{},
			},
			tmdbid:        12345,
			expectedFound: false,
			expectedMALID: 0,
		},
		{
			name: "returns false for nil movie list",
			movies: &AnimeMovies{
				AnimeMovie: nil,
			},
			tmdbid:        12345,
			expectedFound: false,
			expectedMALID: 0,
		},
		{
			name: "finds movie with zero TMDB ID",
			movies: &AnimeMovies{
				AnimeMovie: []AnimeMovie{
					{
						MainTitle: "Test Movie 1",
						TMDBID:    0,
						MALID:     67890,
					},
				},
			},
			tmdbid:        0,
			expectedFound: true,
			expectedMALID: 67890,
		},
		{
			name: "finds movie with zero MAL ID",
			movies: &AnimeMovies{
				AnimeMovie: []AnimeMovie{
					{
						MainTitle: "Test Movie 1",
						TMDBID:    12345,
						MALID:     0,
					},
				},
			},
			tmdbid:        12345,
			expectedFound: true,
			expectedMALID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, movie := tt.movies.CheckMap(tt.tmdbid)
			assert.Equal(t, tt.expectedFound, found)
			if tt.expectedFound {
				require.NotNil(t, movie)
				assert.Equal(t, tt.expectedMALID, movie.MALID)
			} else {
				assert.Nil(t, movie)
			}
		})
	}
}
