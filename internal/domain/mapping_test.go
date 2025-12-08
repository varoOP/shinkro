package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
