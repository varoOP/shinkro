package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

type DB struct {
	handler  *sql.DB
	log      zerolog.Logger
	lock     sync.RWMutex
	squirrel sq.StatementBuilderType
}

func NewDB(dir string, log *zerolog.Logger) *DB {
	db := &DB{
		log: log.With().Str("module", "database").Logger(),
	}

	var (
		err error
		DSN = filepath.Join(dir, "shinkro.db") + "?_pragma=busy_timeout%3d1000"
	)

	db.handler, err = sql.Open("sqlite", DSN)
	if err != nil {
		db.log.Fatal().Err(err).Msg("unable to connect to database")
	}

	if _, err = db.handler.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		if err != nil {
			db.log.Fatal().Err(err).Msg("unable to enable WAL mode")
		}
	}

	return db
}

func (db *DB) Migrate() error {
	db.lock.Lock()
	defer db.lock.Unlock()
	var version int
	if err := db.handler.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return errors.Wrap(err, "failed to query schema version")
	}

	if version == len(migrations) {
		return nil
	} else if version > len(migrations) {
		return errors.Errorf("shinkro (version %d) older than schema (version: %d)", len(migrations), version)
	}

	db.log.Info().Msgf("Beginning database schema upgrade from version %v to version: %v", version, len(migrations))
	tx, err := db.handler.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	if version == 0 {
		if _, err := tx.Exec(schema); err != nil {
			return errors.Wrap(err, "failed to initialize schema")
		}
	} else {
		for i := version; i < len(migrations); i++ {
			db.log.Info().Msgf("Upgrading Database schema to version: %v", i)
			if _, err := tx.Exec(migrations[i]); err != nil {
				return errors.Wrapf(err, "failed to execute migration #%v", i)
			}
		}
	}

	_, err = tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", len(migrations)))
	if err != nil {
		return errors.Wrap(err, "failed to bump schema version")
	}

	db.log.Info().Msgf("Database schema upgraded to version: %v", len(migrations))
	return tx.Commit()
}

func (db *DB) UpdateAnime() {

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

	tx, err := db.handler.Begin()
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

	if _, err = db.handler.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		db.check(err)
	}

	db.log.Trace().Msg("Updated anime in database")
}

func (db *DB) UpdateMalAuth(m map[string]string) {
	const addMalauth = `INSERT OR REPLACE INTO malauth (
		id,
		client_id,
		client_secret,
		access_token
	) values (?, ?, ?, ?)`

	stmt, err := db.handler.Prepare(addMalauth)
	db.check(err)
	defer stmt.Close()

	_, err = stmt.Exec(1, m["client_id"], m["client_secret"], m["access_token"])
	db.check(err)

	if _, err = db.handler.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		db.check(err)
	}
}

func (db *DB) GetMalCreds(ctx context.Context) (map[string]string, error) {
	var (
		client_id     string
		client_secret string
		access_token  string
	)

	sqlstmt := "SELECT client_id, client_secret, access_token from malauth;"

	row := db.handler.QueryRowContext(ctx, sqlstmt)
	err := row.Scan(&client_id, &client_secret, &access_token)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"access_token":  access_token,
	}, nil
}

func (db *DB) Close() error {
	if _, err := db.handler.Exec(`PRAGMA optimize;`); err != nil {
		return errors.Wrap(err, "query planner optimization")
	}

	db.handler.Close()
	return nil
}

func (db *DB) check(err error) {
	if err != nil {
		db.log.Fatal().Err(errors.WithStack(err)).Msg("Database operation failed")
	}
}

func (db *DB) Ping() error {
	return db.handler.Ping()
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.handler.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:      tx,
		handler: db,
	}, nil
}

type Tx struct {
	*sql.Tx
	handler *DB
}
