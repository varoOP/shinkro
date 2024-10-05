package database

const schema = `
CREATE TABLE anime 
(
	mal_id 		INTEGER PRIMARY KEY,
	title 		TEXT,
	en_title 	TEXT,
	anidb_id 	INTEGER,
	tvdb_id 	INTEGER,
	tmdb_id 	INTEGER,
	type 		TEXT,
	releaseDate TEXT
);

CREATE TABLE users
(
    id         INTEGER PRIMARY KEY,
    username   TEXT NOT NULL,
    password   TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (username)
);

CREATE TABLE api_key
(
    name       TEXT,
    key        TEXT PRIMARY KEY,
    scopes     TEXT []   DEFAULT '{}' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE malauth 
(
	id 				INTEGER PRIMARY KEY,
	client_id 		TEXT,
	client_secret 	TEXT,
	access_token 	TEXT
);

CREATE TABLE anime_update
(
	id 				INTEGER PRIMARY KEY,
	mal_id 			INTEGER NOT NULL,
	source_db 		TEXT NOT NULL,
	source_id 		INTEGER NOT NULL,
	episode_num 	INTEGER,
	season_num  	INTEGER,
	time_stamp  	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	list_details 	TEXT,
	list_status  	TEXT,
	plex_id      	INTEGER
	    REFERENCES plex
			ON DELETE SET NULL       
);

CREATE TABLE plex
(
	id 							INTEGER PRIMARY KEY,
	rating              		REAL,
	event  						TEXT NOT NULL,
	source  					TEXT NOT NULL,
	account_title 				TEXT NOT NULL,
	guid_string                 TEXT,
	guids                       TEXT,
	grand_parent_key    		TEXT,
	grand_parent_title  		TEXT,
	metadata_index              INTEGER,
	library_section_title 		TEXT NOT NULL,
	parent_index 				INTEGER,
	title						TEXT,
	type						TEXT NOT NULL,
	time_stamp                  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

var migrations = []string{
	"",
}
