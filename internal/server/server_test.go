package server

import (
	"bytes"
	"database/sql"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/varoOP/shinkuro/internal/config"
)

func TestServer_HandlePlexReq(t *testing.T) {

	t.Skip()

	b, w := createMultipartForm(t, "payload", "testdata/plex_webhook_test.json")
	req := httptest.NewRequest("POST", "/", b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	db := createMockDB(t)
	defer db.Close()

	c := mal.NewClient(http.DefaultClient)
	ac := &AnimeCon{c, db, nil}

	rr := httptest.NewRecorder()

	cfg := &config.Config{
		CustomMap: false,
		User:      "TestUser",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlePlexReq(w, r, ac, cfg)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func createMultipartForm(t *testing.T, fieldname, filepath string) (*bytes.Buffer, *multipart.Writer) {

	body := &bytes.Buffer{}

	w := multipart.NewWriter(body)
	defer w.Close()

	fw, err := w.CreateFormField(fieldname)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(filepath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(fw, f)
	if err != nil {
		t.Fatal(err)
	}

	return body, w

}

func createMockDB(t *testing.T) *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("error creating mock database")
	}

	rows := sqlmock.NewRows([]string{"mal_id"}).AddRow(52305)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	return db
}
