package animedb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/varoOP/shinkuro/pkg/animelist"
	"github.com/varoOP/shinkuro/pkg/manami"
)

func check(err error) {
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
}

func NewDB(DSN string) *sql.DB {

	db, err := sql.Open("sqlite3", DSN)
	check(err)
	return db
}

func UpdateDB(db *sql.DB) {

	err := checkDB(db)
	if err != nil {
		log.Println("Creating DB")
	} else {
		log.Println("Updating DB")
	}

	m := manami.NewManami()
	al := animelist.NewAnimeList()
	am := makeAnimeMap(m, al)

	const scheme = `CREATE TABLE IF NOT EXISTS anime (
		anidb_id INTEGER PRIMARY KEY,
		title TEXT,
		mal_id INTEGER,
		tvdb_id INTEGER
	);`

	_, err = db.Exec(scheme)
	check(err)

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

	log.Println("DB operation complete")

}

func checkDB(db *sql.DB) error {

	sqlstmt := `SELECT * from anime;`

	_, err := db.Query(sqlstmt)
	if err != nil {
		return err
	}

	return nil
}
