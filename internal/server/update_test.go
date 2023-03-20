package server

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/mapping"
	"golang.org/x/oauth2"
)

type have struct {
	data  string
	event string
	cfg   *config.Config
	db    *sql.DB
}

func TestUpdate_TvdbToMal(t *testing.T) {
	tests := []struct {
		name string
		have *AnimeUpdate
		want *AnimeUpdate
	}{
		{
			name: "One Piece",
			have: &AnimeUpdate{
				anime: &mapping.Anime{
					Seasons: []mapping.Seasons{
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
				show: &mapping.Show{
					Season: 21,
					Ep:     162,
				},
				malid: -1,
				start: -1,
			},
			want: &AnimeUpdate{
				malid: 21,
				start: 892,
				show: &mapping.Show{
					Ep: 1053,
				},
			},
		},
		{
			name: "DanMachi",
			have: &AnimeUpdate{
				anime: &mapping.Anime{
					Seasons: []mapping.Seasons{
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
				show: &mapping.Show{
					Season: 4,
					Ep:     13,
				},
				malid: -1,
				start: -1,
			},
			want: &AnimeUpdate{
				malid: 53111,
				start: 12,
				show: &mapping.Show{
					Ep: 2,
				},
			},
		},
		{
			name: "Vinland Saga",
			have: &AnimeUpdate{
				anime: &mapping.Anime{
					Seasons: []mapping.Seasons{
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
				show: &mapping.Show{
					Season: 2,
					Ep:     9,
				},
				malid: -1,
				start: -1,
			},
			want: &AnimeUpdate{
				malid: 49387,
				start: 1,
				show: &mapping.Show{
					Ep: 9,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := tt.have.tvdbtoMal(context.Background())

			if tt.have.malid != tt.want.malid {
				t.Errorf("\nTest: %v\nHave:malid_%v Want:malid_%v", tt.name, tt.have.malid, tt.want.malid)
			}

			if tt.have.start != tt.want.start {
				t.Errorf("\nTest: %v\nHave:start_%v Want:start_%v", tt.name, tt.have.start, tt.want.start)
			}

			if ep != tt.want.show.Ep {
				t.Errorf("\nTest: %v\nHave:ep_%v Want:ep_%v", tt.name, ep, tt.want.show.Ep)
			}
		})
	}
}

func TestUpdate_ServeHTTP(t *testing.T) {

	tests := []struct {
		name string
		have have
		want *mal.AnimeListStatus
	}{
		{
			name: "Tomo-Chan",
			have: have{
				data: `{
				"rating": 8.0,
				"event": "media.rate",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"guid": "com.plexapp.agents.hama://anidb-17494/1/7?lang=en",
					"type": "episode",
					"grandparentTitle": "Tomo-chan wa Onnanoko!"
				}
			}`,
				event: "media.rate",
				cfg: &config.Config{
					CustomMap: "",
					User:      "TestUser",
				},
				db: createMockDB(t, 52305),
			},
			want: &mal.AnimeListStatus{
				Score: 8,
			},
		},
		{
			name: "Isekai Nonbiri Nouka",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"guid": "com.plexapp.agents.hama://anidb-17290/1/9?lang=en",
					"type": "episode",
					"grandparentTitle": "Isekai Nonbiri Nouka"
				}
			}`,
				event: "media.scrobble",
				cfg: &config.Config{
					CustomMap: "",
					User:      "TestUser",
				},
				db: createMockDB(t, 51462),
			},
			want: &mal.AnimeListStatus{
				NumEpisodesWatched: 9,
			},
		},
		{
			name: "Your Name",
			have: have{
				data: `{
				"rating": 8.0,
				"event": "media.rate",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"guid": "net.fribbtastic.coding.plex.myanimelist://32281?lang=en",
					"type": "movie"
				}
			}`,
				event: "media.rate",
				cfg: &config.Config{
					CustomMap: "",
					User:      "TestUser",
				},
				db: createMockDB(t, 0),
			},
			want: &mal.AnimeListStatus{
				Score: 8,
			},
		},
		{
			name: "DanMachi",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"guid": "com.plexapp.agents.hama://tvdb-289882/4/22?lang=en",
					"type": "episode",
					"grandparentTitle": "Dungeon ni Deai o Motomeru no wa Machigatte Iru Darouka: Familia Myth"
				}
			}`,
				event: "media.scrobble",
				cfg: &config.Config{
					CustomMap: "",
					User:      "TestUser",
				},
				db: createMockDB(t, 0),
			},
			want: &mal.AnimeListStatus{
				NumEpisodesWatched: 11,
			},
		},
		{
			name: "Bakemono no Ko",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestUser"
				},
				"Metadata": {
					"guid": "net.fribbtastic.coding.plex.myanimelist://28805?lang=en",
					"type": "movie"
				}
			}`,
				event: "media.scrobble",
				cfg: &config.Config{
					CustomMap: "",
					User:      "TestUser",
				},
				db: createMockDB(t, 0),
			},
			want: &mal.AnimeListStatus{
				NumEpisodesWatched: 1,
			},
		},
	}

	rr := httptest.NewRecorder()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createRequest(t, tt.have.data)
			a := NewAnimeUpdate(tt.have.db, tt.have.cfg)
			a.ServeHTTP(rr, req)
			switch tt.have.event {
			case "media.rate":
				if a.malresp.Score != tt.want.Score {
					t.Errorf("Test:%v Have:%v Want:%v", tt.name, a.malresp.Score, tt.want.Score)
				}
			case "media.scrobble":
				if a.malresp.NumEpisodesWatched != tt.want.NumEpisodesWatched {
					t.Errorf("Test:%v Have:%v Want:%v", tt.name, a.malresp.NumEpisodesWatched, tt.want.NumEpisodesWatched)
				}
			}
		})
	}
}

func createMultipartForm(t *testing.T, data string) (*bytes.Buffer, *multipart.Writer) {

	body := &bytes.Buffer{}

	w := multipart.NewWriter(body)
	defer w.Close()

	fw, err := w.CreateFormField("payload")
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(fw, strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	return body, w

}

func createMalclient(t *testing.T) []string {
	var creds map[string]string
	token := &oauth2.Token{}

	unmarshal(t, "testdata/mal-credentials.json", &creds)
	unmarshal(t, "testdata/token.json", token)

	tt, err := json.Marshal(token)
	if err != nil {
		t.Fatal(err)
	}

	return []string{creds["client-id"], creds["client-secret"], string(tt)}
}

func createMockDB(t *testing.T, malid int) *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("error creating mock database")
	}

	r := createMalclient(t)
	rows := sqlmock.NewRows([]string{"client_id", "client_secret", "access_token"}).AddRow(r[0], r[1], r[2])
	mock.ExpectQuery(`SELECT \* from malauth;`).WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"mal_id"}).AddRow(malid)
	mock.ExpectQuery("SELECT mal_id from anime").WillReturnRows(rows)

	return db
}

func unmarshal(t *testing.T, path string, v any) {
	f, err := os.Open(path)
	if err != nil {
		t.Skip()
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		t.Fatal(err)
	}
}

func createRequest(t *testing.T, data string) *http.Request {

	b, w := createMultipartForm(t, data)
	req := httptest.NewRequest("POST", "/", b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}
