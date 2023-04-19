package domain

import (
	"context"
	"testing"

	"github.com/varoOP/shinkuro/internal/database"
)

func TestUpdateTvdbToMal(t *testing.T) {
	tests := []struct {
		name string
		have *AnimeUpdate
		want *AnimeUpdate
	}{
		{
			name: "One Piece",
			have: &AnimeUpdate{
				Anime: &Anime{
					Seasons: []Seasons{
						{
							Season: 1,
							MalID:  21,
							Start:  0,
						},
						{
							Season: 2,
							MalID:  21,
							Start:  9,
						},
						{
							Season: 3,
							MalID:  21,
							Start:  31,
						},
						{
							Season: 4,
							MalID:  21,
							Start:  48,
						},
						{
							Season: 5,
							MalID:  21,
							Start:  61,
						},
						{
							Season: 6,
							MalID:  21,
							Start:  70,
						},
						{
							Season: 7,
							MalID:  21,
							Start:  92,
						},
						{
							Season: 8,
							MalID:  21,
							Start:  131,
						},
						{
							Season: 9,
							MalID:  21,
							Start:  144,
						},
						{
							Season: 10,
							MalID:  21,
							Start:  196,
						},
						{
							Season: 11,
							MalID:  21,
							Start:  227,
						},
						{
							Season: 12,
							MalID:  21,
							Start:  326,
						},
						{
							Season: 13,
							MalID:  21,
							Start:  382,
						},
						{
							Season: 14,
							MalID:  21,
							Start:  482,
						},
						{
							Season: 15,
							MalID:  21,
							Start:  517,
						},
						{
							Season: 16,
							MalID:  21,
							Start:  579,
						},
						{
							Season: 17,
							MalID:  21,
							Start:  629,
						},
						{
							Season: 18,
							MalID:  21,
							Start:  747,
						},
						{
							Season: 19,
							MalID:  21,
							Start:  780,
						},
						{
							Season: 20,
							MalID:  21,
							Start:  878,
						},
						{
							Season: 21,
							MalID:  21,
							Start:  892,
						},
					},
				},
				Media: &database.Media{
					Season: 21,
					Ep:     162,
				},
				Malid: -1,
				Start: -1,
			},
			want: &AnimeUpdate{
				Malid: 21,
				Start: 892,
				Ep:    1053,
			},
		},
		{
			name: "DanMachi",
			have: &AnimeUpdate{
				Anime: &Anime{
					Seasons: []Seasons{
						{
							Season: 1,
							MalID:  28121,
							Start:  0,
						},
						{
							Season: 2,
							MalID:  37347,
							Start:  0,
						},
						{
							Season: 3,
							MalID:  40454,
							Start:  0,
						},
						{
							Season: 4,
							MalID:  47164,
							Start:  0,
						},
						{
							Season: 4,
							MalID:  53111,
							Start:  12,
						},
					},
				},
				Media: &database.Media{
					Season: 4,
					Ep:     13,
				},
				Malid: -1,
				Start: -1,
			},
			want: &AnimeUpdate{
				Malid: 53111,
				Start: 12,
				Ep:    2,
			},
		},
		{
			name: "Vinland Saga",
			have: &AnimeUpdate{
				Anime: &Anime{
					Seasons: []Seasons{
						{
							Season: 1,
							MalID:  37521,
							Start:  0,
						},
						{
							Season: 2,
							MalID:  49387,
							Start:  0,
						},
					},
				},
				Media: &database.Media{
					Season: 2,
					Ep:     9,
				},
				Malid: -1,
				Start: -1,
			},
			want: &AnimeUpdate{
				Malid: 49387,
				Start: 1,
				Ep:    9,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.have.tvdbtoMal(context.Background())
			if err != nil {
				t.Error(err)
			}

			if tt.have.Malid != tt.want.Malid {
				t.Errorf("\nTest: %v\nHave:malid_%v Want:malid_%v", tt.name, tt.have.Malid, tt.want.Malid)
			}

			if tt.have.Start != tt.want.Start {
				t.Errorf("\nTest: %v\nHave:start_%v Want:start_%v", tt.name, tt.have.Start, tt.want.Start)
			}

			if tt.have.Ep != tt.want.Ep {
				t.Errorf("\nTest: %v\nHave:ep_%v Want:ep_%v", tt.name, tt.have.Ep, tt.want.Ep)
			}
		})
	}
}
