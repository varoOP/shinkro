package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/pkg/plex"
	"golang.org/x/oauth2"
)

type have struct {
	data  string
	event string
	cfg   *domain.Config
	db    *database.DB
}

const (
	scrobbleEvent = "media.scrobble"
	rateEvent     = "media.rate"
)

func TestPlex(t *testing.T) {

	tests := []struct {
		name string
		have have
	}{
		{
			name: "HAMA_Episode_DB_Rate_1",
			have: have{
				data: `{
				"rating": 8.0,
				"event": "media.rate",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "com.plexapp.agents.hama://anidb-17494/1/7?lang=en",
					"type": "episode",
					"parentIndex": 1,
					"index": 7,
					"grandparentTitle": "Tomo-chan wa Onnanoko!"
				}
			}`,
				event: rateEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 52305),
			},
		},
		{
			name: "HAMA_Episode_DB_Scrobble_1",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "com.plexapp.agents.hama://anidb-17290/1/9?lang=en",
					"type": "episode",
					"parentIndex": 1,
					"index": 9,
					"grandparentTitle": "Isekai Nonbiri Nouka"
				}
			}`,
				event: scrobbleEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 51462),
			},
		},
		{
			name: "HAMA_Episode_Mapping_Scrobble_1",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "com.plexapp.agents.hama://tvdb-289882/4/22?lang=en",
					"type": "episode",
					"parentIndex": 4,
					"index": 22,
					"grandparentTitle": "Dungeon ni Deai o Motomeru no wa Machigatte Iru Darouka: Familia Myth"
				}
			}`,
				event: scrobbleEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 0),
			},
		},
		{
			name: "HAMA_Episode_Mapping_Scrobble_2",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "com.plexapp.agents.hama://tvdb-316842/0/38?lang=en",
					"type": "episode",
					"parentIndex": 0,
					"index": 38,
					"grandparentTitle": "Mahou Tsukai no Yome"
				}
			}`,
				event: scrobbleEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 0),
			},
		},
		{
			name: "MAL_Movie_Scrobble_1",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "net.fribbtastic.coding.plex.myanimelist://28805?lang=en",
					"type": "movie",
					"parentIndex": 1,
					"index": 1,
					"title":"Mal_Movie"
				}
			}`,
				event: scrobbleEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 0),
			},
		},
		{
			name: "MAL_Movie_Rate_1",
			have: have{
				data: `{
				"rating": 8.0,
				"event": "media.rate",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "net.fribbtastic.coding.plex.myanimelist://32281?lang=en",
					"type": "movie",
					"parentIndex": 1,
					"index": 1,
					"title": "Mal_Movie"
				}
			}`,
				event: rateEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 0),
			},
		},
		{
			name: "MAL_Episode_Scrobble_1",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "net.fribbtastic.coding.plex.myanimelist://52173/1/5?lang=en",
					"type": "episode",
					"parentIndex": 1,
					"index": 5,
					"grandparentTitle": "Mal_Episode"
				}
			}`,
				event: scrobbleEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 0),
			},
		},
		{
			name: "MAL_Episode_Scrobble_2",
			have: have{
				data: `{
				"event": "media.scrobble",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "net.fribbtastic.coding.plex.myanimelist://47160/2/2?lang=en",
					"type": "episode",
					"parentIndex": 2,
					"index": 2,
					"grandparentTitle": "Goblin Slayer II"
				}
			}`,
				event: scrobbleEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 0),
			},
		},
		{
			name: "MAL_Episode_Rate_1",
			have: have{
				data: `{
				"rating": 7.0,
				"event": "media.rate",
				"Account": {
					"title": "TestPlexUser"
				},
				"Metadata": {
					"guid": "net.fribbtastic.coding.plex.myanimelist://52305/1/7?lang=en",
					"type": "episode",
					"parentIndex": 1,
					"index": 7,
					"grandparentTitle": "Mal_Episode"
				}
			}`,
				event: rateEvent,
				cfg: &domain.Config{
					PlexUser: "TestPlexUser",
				},
				db: createMockDB(t, 0),
			},
		},
	}

	rr := httptest.NewRecorder()
	log := zerolog.New(os.Stdout).With().Logger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createRequest(t, tt.have.data)
			ServeHTTP := Plex(tt.have.db, tt.have.cfg, &log, &domain.Notification{})
			ServeHTTP(rr, req)
			if rr.Result().StatusCode != 204 {
				t.Errorf("%s test failed", tt.name)
			}
		})
	}
}

// func createMultipartForm(t *testing.T, data string) (*bytes.Buffer, *multipart.Writer) {

// 	body := &bytes.Buffer{}

// 	w := multipart.NewWriter(body)
// 	defer w.Close()

// 	fw, err := w.CreateFormField("payload")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = io.Copy(fw, strings.NewReader(data))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	return body, w

// }

func createPlexPayload(data string) (*plex.PlexWebhook, error) {
	p, err := plex.NewPlexWebhook([]byte(data))
	if err != nil {
		return nil, err
	}

	return p, nil
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

func createMockDB(t *testing.T, malid int) *database.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("error creating mock database")
	}

	r := createMalclient(t)
	rows := sqlmock.NewRows([]string{"client_id", "client_secret", "access_token"}).AddRow(r[0], r[1], r[2])
	mock.ExpectQuery(`SELECT client_id, client_secret, access_token from malauth;`).WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"mal_id"}).AddRow(malid)
	mock.ExpectQuery("SELECT mal_id from anime").WillReturnRows(rows)

	return &database.DB{
		Handler: db,
	}
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
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	p, err := createPlexPayload(data)
	if err != nil {
		t.Errorf("failed to create plex payload. err: %v", err)
	}

	_, agent := isMetadataAgent(p)

	ctx := req.Context()
	ctx = context.WithValue(ctx, domain.PlexPayload, p)
	ctx = context.WithValue(ctx, domain.Agent, agent)
	return req.WithContext(ctx)
}
