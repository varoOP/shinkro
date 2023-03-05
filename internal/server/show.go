package server

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

type Show struct {
	idSource string
	id       int
	season   int
	ep       int
}

func NewShow(ctx context.Context, guid string) (*Show, error) {

	var err error

	s := &Show{}

	r := regexp.MustCompile(`//(.* ?)-(\d+ ?)/?(\d+ ?)?/?(\d+ ?)?`)

	if !r.MatchString(guid) {
		return s, fmt.Errorf("unable to parse GUID: %v", guid)
	}

	m := r.FindStringSubmatch(guid)

	s.idSource = m[1]

	s.id, err = strconv.Atoi(m[2])
	if err != nil {
		return s, err
	}

	s.season, err = strconv.Atoi(m[3])
	if err != nil {
		return s, err
	}

	s.ep, err = strconv.Atoi(m[4])
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s *Show) GetMalID(ctx context.Context, db *sql.DB) (int, error) {

	var malid int
	sqlstmt := fmt.Sprintf("SELECT mal_id from anime where %v_id=?;", s.idSource)

	row := db.QueryRow(sqlstmt, s.id)
	err := row.Scan(&malid)
	if err != nil {
		return -1, fmt.Errorf("mal_id of %v %v not found in DB", s.idSource, s.id)
	}

	return malid, nil
}
