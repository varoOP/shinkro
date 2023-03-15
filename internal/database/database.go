package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/varoOP/shinkuro/pkg/animelist"
	"github.com/varoOP/shinkuro/pkg/manami"
)

func NewDB(DSN string) *sql.DB {

	db, err := sql.Open("sqlite3", DSN)
	check(err)
	return db
}

func CreateDB(db *sql.DB) {
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

	_, err := db.Exec(scheme)
	check(err)
}

func UpdateAnime(db *sql.DB) {

	check(checkDB(db))
	log.Println("Updating DB")

	m := manami.NewManami()
	al := animelist.NewAnimeList()
	am := makeAnimeMap(m, al)

	var addAnime = `INSERT OR REPLACE INTO anime (
		anidb_id,
		title,
		mal_id,
		tvdb_id
	) values (?, ?, ?, ?)`

	stmt, err := db.Prepare(addAnime)
	check(err)

	defer stmt.Close()

	for _, anime := range am.Anime {
		_, err := stmt.Exec(anime.AnidbID, anime.Title, anime.MalID, anime.TvdbID)
		check(err)

	}

	m, al, am = nil, nil, nil

	log.Println("DB operation complete")

}

func UpdateMalAuth(m map[string]string, db *sql.DB) {
	var addMalauth = `INSERT OR REPLACE INTO malauth (
		client_id,
		client_secret,
		access_token
	) values (?, ?, ?)`

	stmt, err := db.Prepare(addMalauth)
	check(err)
	defer stmt.Close()

	_, err = stmt.Exec(m["client_id"], m["client_secret"], m["access_token"])
	check(err)
}

func GetMalCreds(db *sql.DB) map[string]string {
	var (
		client_id     string
		client_secret string
		access_token  string
	)

	sqlstmt := "SELECT * from malauth;"

	row := db.QueryRow(sqlstmt)
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

func checkDB(db *sql.DB) error {

	sqlstmt := `SELECT * from anime;`

	_, err := db.Query(sqlstmt)
	if err != nil {
		return err
	}

	return nil
}

func check(err error) {
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
}
