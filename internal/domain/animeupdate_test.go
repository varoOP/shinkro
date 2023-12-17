package domain

// import (
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	// "github.com/varoOP/shinkro/internal/database"
// )

// func TestUpdateTvdbToMal(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		have  *AnimeUpdate
// 		want1 int
// 		want2 int
// 		want3 int
// 	}{
// 		{
// 			name: "One Piece",
// 			have: &AnimeUpdate{
// 				Anime: &Anime{
// 					Malid:      21,
// 					UseMapping: true,
// 					TvdbSeason: 21,
// 					Start:      892,
// 				},
// 				Media: &database.Media{
// 					Season: 21,
// 					Ep:     162,
// 				},
// 				Malid: -1,
// 				Start: -1,
// 			},
// 			want1: 21,
// 			want2: 892,
// 			want3: 1053,
// 		},
// 		{
// 			name: "DanMachi",
// 			have: &AnimeUpdate{
// 				Anime: &Anime{
// 					Malid:      53111,
// 					UseMapping: false,
// 					TvdbSeason: 4,
// 					Start:      12,
// 				},
// 				Media: &database.Media{
// 					Season: 4,
// 					Ep:     13,
// 				},
// 				Malid: -1,
// 				Start: -1,
// 			},
// 			want1: 53111,
// 			want2: 12,
// 			want3: 2,
// 		},
// 		{
// 			name: "Vinland Saga",
// 			have: &AnimeUpdate{
// 				Anime: &Anime{
// 					Malid:      49387,
// 					TvdbSeason: 2,
// 					Start:      0,
// 					UseMapping: false,
// 				},
// 				Media: &database.Media{
// 					Season: 2,
// 					Ep:     9,
// 				},
// 				Malid: -1,
// 				Start: -1,
// 			},
// 			want1: 49387,
// 			want2: 1,
// 			want3: 9,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.have.tvdbtoMal(context.Background())
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			assert.Equal(t, tt.have.Malid, tt.want1)
// 			assert.Equal(t, tt.have.Start, tt.want2)
// 			assert.Equal(t, tt.have.Ep, tt.want3)
// 		})
// 	}
// }
