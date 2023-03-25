package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/varoOP/shinkuro/pkg/animelist"
	"github.com/varoOP/shinkuro/pkg/manami"
	_ "modernc.org/sqlite"
)

type DB struct {
	Handler *sql.DB
}

func NewDB(DSN string) *DB {
	db := &DB{}
	var err error

	db.Handler, err = sql.Open("sqlite", DSN)
	check(err)

	if _, err = db.Handler.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		check(err)
	}

	return db
}

func (db *DB) CreateDB() {
	const scheme = `CREATE TABLE IF NOT EXISTS anime (
		anidb_id INTEGER PRIMARY KEY,
		title TEXT,
		mal_id INTEGER,
		tvdb_id INTEGER
	);
	CREATE TABLE IF NOT EXISTS malauth (
		client_id TEXT PRIMARY KEY,
		client_secret TEXT,
		access_token TEXT
	);`

	_, err := db.Handler.Exec(scheme)
	check(err)
}

func (db *DB) UpdateAnime() {

	check(db.checkDB())
	log.Println("Updating DB")

	m := manami.NewManami()
	al := animelist.NewAnimeList()
	am := makeAnimeMap(m, al)

	const addAnime = `INSERT OR REPLACE INTO anime (
		anidb_id,
		title,
		mal_id,
		tvdb_id
	) values (?, ?, ?, ?)`

	tx, err := db.Handler.Begin()
	check(err)

	defer tx.Rollback()

	stmt, err := tx.Prepare(addAnime)
	check(err)

	defer stmt.Close()

	for _, anime := range am.Anime {
		_, err := stmt.Exec(anime.AnidbID, anime.Title, anime.MalID, anime.TvdbID)
		check(err)
	}

	if err = tx.Commit(); err != nil {
		check(err)
	}

	if _, err = db.Handler.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		check(err)
	}

	log.Println("DB operation complete")
}

func (db *DB) UpdateMalAuth(m map[string]string) {
	const addMalauth = `INSERT OR REPLACE INTO malauth (
		client_id,
		client_secret,
		access_token
	) values (?, ?, ?)`

	stmt, err := db.Handler.Prepare(addMalauth)
	check(err)
	defer stmt.Close()

	_, err = stmt.Exec(m["client_id"], m["client_secret"], m["access_token"])
	check(err)

	if _, err = db.Handler.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		check(err)
	}
}

func (db *DB) GetMalCreds(ctx context.Context) map[string]string {
	var (
		client_id     string
		client_secret string
		access_token  string
	)

	sqlstmt := "SELECT * from malauth;"

	row := db.Handler.QueryRowContext(ctx, sqlstmt)
	err := row.Scan(&client_id, &client_secret, &access_token)
	if err != nil {
		check(err)
	}

	return map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"access_token":  access_token,
	}
}

func (db *DB) checkDB() error {

	sqlstmt := `SELECT * from anime;`

	_, err := db.Handler.Query(sqlstmt)
	if err != nil {
		return errors.New("shinkuro.db not found. Run shinkuro malauth before running shinkuro")
	}
	return nil
}

func (db *DB) Close() {
	db.Handler.Close()
}

func check(err error) {
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
}
