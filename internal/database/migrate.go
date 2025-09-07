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
    admin      BOOLEAN DEFAULT false NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (username)
);

CREATE TABLE api_key
(
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT,
    key        TEXT PRIMARY KEY,
    scopes     TEXT []   DEFAULT '{}' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE malauth 
(
	id 					INTEGER PRIMARY KEY,
	user_id				INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	client_id 			TEXT,
	client_secret 		TEXT,
	access_token 		BLOB,
	token_iv 			BLOB,
	created_at 			TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at 			TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE anime_update
(
	id 				INTEGER PRIMARY KEY,
	user_id			INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
			ON DELETE CASCADE       
);

CREATE TABLE plex_payload
(
	id 							INTEGER PRIMARY KEY,
	user_id						INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
	user_id						INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	host						TEXT NOT NULL,
	port						INTEGER,
	tls 						BOOLEAN,
	tls_skip_verify				BOOLEAN,
	token						BLOB,
	token_iv					BLOB,
	username					TEXT NOT NULL,
	anime_libraries             TEXT []   DEFAULT '{}' NOT NULL,
	plex_client_enabled			BOOLEAN DEFAULT false NOT NULL,
	client_id				    TEXT,
	time_stamp                  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE plex_status
(
    	id 							INTEGER PRIMARY KEY,
    	title 							TEXT NOT NULL,
    	event 						TEXT NOT NULL,
    	success 			        BOOLEAN NOT NULL,
    	error_msg 				TEXT,
    	time_stamp 			TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	plex_id 					INTEGER NOT NULL
    	  REFERENCES plex_payload
		     ON DELETE CASCADE
);

CREATE TABLE mapping_settings
(
	id						        INTEGER PRIMARY KEY,
	user_id						INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	tvdb_enabled			BOOLEAN,
	tmdb_enabled		BOOLEAN,
	tvdb_path				TEXT,
	tmdb_path				TEXT,
	time_stamp             TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notification
(
	id         INTEGER PRIMARY KEY,
	user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	name       TEXT,
	type       TEXT,
	enabled    BOOLEAN,
	events     TEXT []   DEFAULT '{}' NOT NULL,
	token      TEXT,
	api_key    TEXT,
	webhook    TEXT,
	title      TEXT,
	icon       TEXT,
	host       TEXT,
	username   TEXT,
	password   TEXT,
	channel    TEXT,
	rooms      TEXT,
	targets    TEXT,
	devices    TEXT,
	topic      TEXT,
	priority   INTEGER DEFAULT 0,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

var migrations = []string{
	"",
	`ALTER TABLE users ADD COLUMN admin BOOLEAN DEFAULT false NOT NULL;
	 ALTER TABLE malauth ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
	 ALTER TABLE plex_settings ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
	 UPDATE malauth SET user_id = 1 WHERE user_id IS NULL;
	 UPDATE plex_settings SET user_id = 1 WHERE user_id IS NULL;
	 ALTER TABLE api_key ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
	 UPDATE api_key SET user_id = 1 WHERE user_id IS NULL;
	 ALTER TABLE anime_update ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
	 UPDATE anime_update SET user_id = 1 WHERE user_id IS NULL;
	 ALTER TABLE plex_payload ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
	 UPDATE plex_payload SET user_id = 1 WHERE user_id IS NULL;
	 ALTER TABLE notification ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
	 UPDATE notification SET user_id = 1 WHERE user_id IS NULL;
	 ALTER TABLE mapping_settings ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
	 UPDATE mapping_settings SET user_id = 1 WHERE user_id IS NULL;
	 UPDATE users SET admin = true WHERE id IN (SELECT id FROM users LIMIT 1) AND (SELECT COUNT(*) FROM users) = 1;`,
}
