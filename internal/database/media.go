package database

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/varoOP/shinkro/pkg/plex"
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

var AgentRegExMap = map[string]string{
	"hama": `//(.* ?)-(\d+ ?)`,
	"mal":  `.(m.*)://(\d+ ?)`,
}

func NewMedia(pw *plex.PlexWebhook, agent string, pc *plex.PlexClient, usePlex bool) (*Media, error) {
	var (
		idSource  string
		title     string = pw.Metadata.GrandparentTitle
		mediaType string = pw.Metadata.Type
		id        int
		season    int = pw.Metadata.ParentIndex
		ep        int = pw.Metadata.Index
		err       error
	)

	switch agent {
	case "mal", "hama":
		idSource, id, err = hamaMALAgent(pw.Metadata.GUID.GUID, agent)
		if err != nil {
			return nil, err
		}

	case "plex":
		if !usePlex {
			return nil, errors.New("plex token not provided or invalid")
		}

		guid := pw.Metadata.GUID
		if mediaType == "episode" {
			g, err := GetShowID(pc, pw.Metadata.GrandparentKey)
			if err != nil {
				return nil, err
			}

			guid = *g
		}

		idSource, id, err = plexAgent(guid, mediaType)
		if err != nil {
			return nil, err
		}
	}

	if idSource == "myanimelist" {
		idSource = "mal"
	}

	if mediaType == "movie" {
		season = 1
		ep = 1
		title = pw.Metadata.Title
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
	case "mal":
		malid = m.Id

	default:
		sqlstmt := fmt.Sprintf("SELECT mal_id from anime where %v_id=?;", m.IdSource)
		row := db.handler.QueryRowContext(ctx, sqlstmt, m.Id)
		err := row.Scan(&malid)
		if err != nil {
			return -1, errors.Errorf("mal_id of %v (%v:%v) not found in database, add to custom map", m.Title, m.IdSource, m.Id)
		}
	}

	return malid, nil
}

func (m *Media) ConvertToTVDB(ctx context.Context, db *DB) {
	if m.IdSource == "anidb" && m.Season > 1 {
		var tvdbid int
		sqlstmt := "SELECT tvdb_id from anime where anidb_id=?;"
		row := db.handler.QueryRowContext(ctx, sqlstmt, m.Id)
		err := row.Scan(&tvdbid)
		if err != nil {
			return
		}

		m.IdSource = "tvdb"
		m.Id = tvdbid
	}
}

func hamaMALAgent(guid, agent string) (string, int, error) {
	r := regexp.MustCompile(AgentRegExMap[agent])
	if !r.MatchString(guid) {
		return "", -1, errors.Errorf("unable to parse GUID: %v", guid)
	}

	mm := r.FindStringSubmatch(guid)
	source := mm[1]
	id, err := strconv.Atoi(mm[2])
	if err != nil {
		return "", -1, errors.Wrap(err, "conversion of id failed")
	}

	return source, id, nil
}

func plexAgent(guid plex.GUID, mediaType string) (string, int, error) {
	for _, gid := range guid.GUIDS {
		dbid := strings.Split(gid.ID, "://")
		if (mediaType == "episode" && dbid[0] == "tvdb") || (mediaType == "movie" && dbid[0] == "tmdb") {
			id, err := strconv.Atoi(dbid[1])
			if err != nil {
				return "", -1, errors.Wrap(err, "id conversion failed")
			}

			return dbid[0], id, nil
		}
	}

	return "", -1, errors.New("no supported db found")
}

func GetShowID(p *plex.PlexClient, key string) (*plex.GUID, error) {
	guid, err := p.GetShowID(key)
	if err != nil {
		return nil, err
	}

	return guid, nil
}
