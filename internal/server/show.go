package server

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type Show struct {
	IdSource string
	Id       int
	Ep       Episode
}

type Episode struct {
	Season int
	No     int
}

func NewShow(guid string) *Show {

	var err error

	s := &Show{}

	r := regexp.MustCompile(`^com.plexapp.agents.hama://(tvdb|anidb)-([0-9]+)(?:/([0-9]+)/([0-9]+))?\?lang=\w{2}$`)
	matches := r.FindStringSubmatch(guid)

	s.IdSource = matches[1]

	s.Id, err = strconv.Atoi(matches[2])
	if err != nil {
		log.Fatalf("error converting anime id from string to inr: %v", err)
	}

	s.Ep.Season, err = strconv.Atoi(matches[3])
	if err != nil {
		s.Ep.Season = 0
	}

	s.Ep.No, err = strconv.Atoi(matches[4])
	if err != nil {
		s.Ep.No = 0
	}

	return s
}

func (s *Show) GetMalID(db *sql.DB) int {

	var malid int
	sqlstmt := fmt.Sprintf("SELECT mal_id from anime where %v_id=?;", s.IdSource)

	row := db.QueryRow(sqlstmt, s.Id)
	err := row.Scan(&malid)
	if err != nil {
		return 0
	}

	return malid
}
