package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMappingCheckMap(t *testing.T) {
	tests := []struct {
		name       string
		animeMap   *AnimeTVDBMap
		tvdbid     int
		tvdbseason int
		ep         int
		want1      bool
		want2      int //malid
		want3      int //start
	}{
		{
			name: "DanMachi S4-1",
			animeMap: &AnimeTVDBMap{
				Anime: []Anime{
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
			animeMap: &AnimeTVDBMap{
				Anime: []Anime{
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
			name: "One Piece-1",
			animeMap: &AnimeTVDBMap{
				Anime: []Anime{
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
			animeMap: &AnimeTVDBMap{
				Anime: []Anime{
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
