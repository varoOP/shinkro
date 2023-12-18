package database

const schema = `
CREATE TABLE IF NOT EXISTS anime 
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

CREATE TABLE IF NOT EXISTS malauth 
(
	id 				INTEGER PRIMARY KEY,
	client_id 		TEXT,
	client_secret 	TEXT,
	access_token 	TEXT
);

CREATE TABLE anime_update
(
	mal_id 					INTEGER PRIMARY KEY,
	title 					TEXT NOT NULL,
	status 					TEXT NOT NULL,
	score 					INTEGER,
	num_episodes_watched 	INTEGER,
	is_rewatching 			BOOLEAN,
	updated_at 				TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	priority 				INTEGER,
	num_times_rewatched 	INTEGER,
	rewatch_value 			INTEGER,
	tags 					TEXT []	DEFAULT '{}',
	comments 				TEXT,
	start_date 				TEXT,
	finish_date 			TEXT,
	mal_url					TEXT
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
