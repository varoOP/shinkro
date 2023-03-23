package database

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
)

type Show struct {
	IdSource string
	Id       int
	Season   int
	Ep       int
}

func NewShow(guid string) (*Show, error) {

	r := regexp.MustCompile(`//(.* ?)-(\d+ ?)/?(\d+ ?)?/?(\d+ ?)?`)

	if !r.MatchString(guid) {
		return nil, fmt.Errorf("unable to parse GUID: %v", guid)
	}

	m := r.FindStringSubmatch(guid)

	idSource := m[1]

	id, err := strconv.Atoi(m[2])
	if err != nil {
		return nil, err
	}

	season, err := strconv.Atoi(m[3])
	if err != nil {
		return nil, err
	}

	ep, err := strconv.Atoi(m[4])
	if err != nil {
		return nil, err
	}

	return &Show{idSource, id, season, ep}, nil
}

func (s *Show) GetMalID(ctx context.Context, db *DB) (int, error) {

	var malid int
	sqlstmt := fmt.Sprintf("SELECT mal_id from anime where %v_id=?;", s.IdSource)

	row := db.Handler.QueryRow(sqlstmt, s.Id)
	err := row.Scan(&malid)
	if err != nil {
		return -1, fmt.Errorf("mal_id of %v %v not found in DB", s.IdSource, s.Id)
	}

	return malid, nil
}
