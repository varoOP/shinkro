package database

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Media struct {
	Type     string
	Title    string
	Agent    string
	IdSource string
	Id       int
	Season   int
	Ep       int
}

func NewMedia(guid, mediaType, title string) (*Media, error) {
	var (
		agent    string
		idSource string
		id       int
		season   int
		ep       int
		err      error
	)

	r := regexp.MustCompile(`//(.* ?)-(\d+ ?)/?(\d+ ?)?/?(\d+ ?)?`)
	agent = "hama"

	if strings.Contains(guid, "net.fribbtastic.coding.plex.myanimelist") {
		r = regexp.MustCompile(`(myanimelist)://(\d+ ?)/?(\d+ ?)?/?(\d+ ?)?`)
		agent = "mal"
	}

	if !r.MatchString(guid) {
		return nil, fmt.Errorf("unable to parse GUID: %v", guid)
	}

	mm := r.FindStringSubmatch(guid)

	switch mediaType {
	case "episode":
		idSource = mm[1]
		id, err = strconv.Atoi(mm[2])
		if err != nil {
			return nil, err
		}

		season, err = strconv.Atoi(mm[3])
		if err != nil {
			return nil, err
		}

		ep, err = strconv.Atoi(mm[4])
		if err != nil {
			return nil, err
		}

	case "movie":
		if agent == "hama" {
			return nil, errors.New("hama agent for movies not supported")
		}

		idSource = "mal"
		id, err = strconv.Atoi(mm[2])
		if err != nil {
			return nil, err
		}

		season = 1
		ep = 1

	default:
		return nil, fmt.Errorf("%v media type not supported", mediaType)
	}

	return &Media{
		Type:     mediaType,
		Title:    title,
		Agent:    agent,
		IdSource: idSource,
		Id:       id,
		Season:   season,
		Ep:       ep,
	}, nil
}

func (m *Media) GetMalID(ctx context.Context, db *DB) (int, error) {
	var malid int
	switch m.Agent {
	case "hama":
		sqlstmt := fmt.Sprintf("SELECT mal_id from anime where %v_id=?;", m.IdSource)
		row := db.Handler.QueryRowContext(ctx, sqlstmt, m.Id)
		err := row.Scan(&malid)
		if err != nil {
			return -1, fmt.Errorf("mal_id of %v (%v:%v) not found in database",m.Title, m.IdSource, m.Id)
		}
	case "mal":
		malid = m.Id
	}

	return malid, nil
}
