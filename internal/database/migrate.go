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
	releaseDate TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
	id 					INTEGER PRIMARY KEY,
	client_id 			TEXT,
	client_secret 		TEXT,
	access_token 		TEXT,
	created_at 			TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at 			TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
	    REFERENCES plex_payload
			ON DELETE SET NULL       
);

CREATE TABLE plex_payload
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

CREATE TABLE plex_settings
(
	id 							INTEGER PRIMARY KEY,
	host						TEXT NOT NULL,
	port						INTEGER,
	tls 						BOOLEAN,
	tls_skip_verify				BOOLEAN,
	token						TEXT,
	username					TEXT NOT NULL,
	anime_libraries             TEXT []   DEFAULT '{}' NOT NULL,
	plex_client_enabled			BOOLEAN DEFAULT false NOT NULL,
	client_id				    TEXT,
	time_stamp                  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE mapping_settings
(
	id						INTEGER PRIMARY KEY,
	tvdb_enabled			BOOLEAN,
	tmdb_enabled			BOOLEAN,
	tvdb_path				TEXT,
	tmdb_path				TEXT,
	created_at TIMESTAMP 	DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP 	DEFAULT CURRENT_TIMESTAMP
);
`

var migrations = []string{
	"",
}
