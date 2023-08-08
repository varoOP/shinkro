package database

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"

	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

type DB struct {
	Handler *sql.DB
	log     zerolog.Logger
}

func NewDB(dir string, log *zerolog.Logger) *DB {
	db := &DB{
		log: log.With().Str("module", "database").Logger(),
	}

	var (
		err error
		DSN = filepath.Join(dir, "shinkro.db") + "?_pragma=busy_timeout%3d1000"
	)

	db.Handler, err = sql.Open("sqlite", DSN)
	if err != nil {
		db.log.Fatal().Err(err).Msg("unable to connect to database")
	}

	if _, err = db.Handler.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		if err != nil {
			db.log.Fatal().Err(err).Msg("unable to enable WAL mode")
		}
	}

	return db
}

func (db *DB) CreateDB() {
	const scheme = `CREATE TABLE IF NOT EXISTS anime (
		mal_id INTEGER PRIMARY KEY,
		title TEXT,
		en_title TEXT,
		anidb_id INTEGER,
		tvdb_id INTEGER,
		tmdb_id INTEGER,
		type TEXT,
		releaseDate TEXT
	);
	CREATE TABLE IF NOT EXISTS malauth (
		client_id TEXT PRIMARY KEY,
		client_secret TEXT,
		access_token TEXT
	);`

	_, err := db.Handler.Exec(scheme)
	db.check(err)
}

func (db *DB) UpdateAnime() {

	db.check(db.checkDBForCreds())
	db.log.Trace().Msg("Updating anime in database")
	a, err := getAnime()
	if err != nil {
		db.log.Error().Err(err).Msg("Unable to update anime in database")
		return
	}

	const addAnime = `INSERT OR REPLACE INTO anime (
		mal_id,
		title,
		en_title,
		anidb_id,
		tvdb_id,
		tmdb_id,
		type,
		releaseDate
	) values (?, ?, ?, ?, ?, ?, ?, ?)`

	tx, err := db.Handler.Begin()
	db.check(err)

	defer tx.Rollback()

	stmt, err := tx.Prepare(addAnime)
	db.check(err)

	defer stmt.Close()

	for _, anime := range a {
		_, err := stmt.Exec(anime.MalID, anime.MainTitle, anime.EnglishTitle, anime.AnidbID, anime.TvdbID, anime.TmdbID, anime.Type, anime.ReleaseDate)
		db.check(err)
	}

	if err = tx.Commit(); err != nil {
		db.check(err)
	}

	if _, err = db.Handler.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		db.check(err)
	}

	db.log.Trace().Msg("Updated anime in database")
}

func (db *DB) UpdateMalAuth(m map[string]string) {
	const addMalauth = `INSERT OR REPLACE INTO malauth (
		client_id,
		client_secret,
		access_token
	) values (?, ?, ?)`

	stmt, err := db.Handler.Prepare(addMalauth)
	db.check(err)
	defer stmt.Close()

	_, err = stmt.Exec(m["client_id"], m["client_secret"], m["access_token"])
	db.check(err)

	if _, err = db.Handler.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		db.check(err)
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
		db.log.Fatal().Str("Possible fix", "Run the command: shinkro malauth").Err(errors.New("unable to get mal credentials")).Msg("Database operation failed")
	}

	return map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"access_token":  access_token,
	}
}

func (db *DB) checkDBForCreds() error {
	db.GetMalCreds(context.Background())
	return nil
}

func (db *DB) Close() {
	db.Handler.Close()
}

func (db *DB) check(err error) {
	if err != nil {
		db.log.Fatal().Err(err).Msg("Database operation failed")
	}
}
