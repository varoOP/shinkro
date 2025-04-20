package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type PlexRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewPlexRepo(log zerolog.Logger, db *DB) domain.PlexRepo {
	return &PlexRepo{
		log: log.With().Str("repo", "plex").Logger(),
		db:  db,
	}
}

func (repo *PlexRepo) Store(ctx context.Context, r *domain.Plex) error {

	guids, err := json.Marshal(r.Metadata.GUID.GUIDS)
	if err != nil {
		return errors.Wrap(err, "failed to marshal GUIDs")
	}

	queryBuilder := repo.db.squirrel.
		Insert("plex_payload").
		Columns("rating", "event", "source", "account_title", "guid_string", "guids", "grand_parent_key", "grand_parent_title", "metadata_index", "library_section_title", "parent_index", "title", "type", "time_stamp").
		Values(r.Rating, r.Event, r.Source, r.Account.Title, r.Metadata.GUID.GUID, string(guids), r.Metadata.GrandparentKey, r.Metadata.GrandparentTitle, r.Metadata.Index, r.Metadata.LibrarySectionTitle, r.Metadata.ParentIndex, r.Metadata.Title, r.Metadata.Type, r.TimeStamp.Format(time.RFC3339)).
		Suffix("RETURNING id").RunWith(repo.db.handler)

	var retID int64

	if err := queryBuilder.QueryRowContext(ctx).Scan(&retID); err != nil {
		repo.log.Debug().Err(err).Msg("error executing query")
		return errors.Wrap(err, "error executing query")
	}

	r.ID = retID
	repo.log.Debug().Msgf("plex.store: %+v", r)
	return nil
}

func (repo *PlexRepo) FindAll(ctx context.Context) ([]*domain.Plex, error) {
	return nil, nil
}

func (repo *PlexRepo) Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error) {
	queryBuilder := repo.db.squirrel.
		Select("p.id", "p.rating", "p.event", "p.source", "p.account_title", "p.guid_string", "p.guids", "p.grand_parent_key", "p.grand_parent_title", "p.metadata_index", "p.library_section_title", "p.parent_index", "p.title", "p.type", "p.time_stamp").
		From("plex_payload p").
		OrderBy("p.id DESC").
		Where(sq.Eq{"p.id": req.Id})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "plex.get").Msgf("query: '%s', args: '%v'", query, args)

	row := repo.db.handler.QueryRowContext(ctx, query, args...)

	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error rows find plex")
	}

	var plex domain.Plex
	var rating sql.NullFloat64
	var index, parent_index sql.NullInt32
	var grandParentKey, grandParentTitle, title, guids, guid_string sql.NullString

	if err := row.Scan(&plex.ID, &rating, &plex.Event, &plex.Source, &plex.Account.Title, &guid_string, &guids, &grandParentKey, &grandParentTitle, &index, &plex.Metadata.LibrarySectionTitle, &parent_index, &title, &plex.Metadata.Type, &plex.TimeStamp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, errors.Wrap(err, "error scanning row")
	}

	err = json.Unmarshal([]byte(guids.String), &plex.Metadata.GUID.GUIDS)
	if err != nil {
		return nil, err
	}

	plex.Metadata.GUID.GUID = guid_string.String
	plex.Rating = float32(rating.Float64)
	plex.Metadata.GrandparentKey = grandParentKey.String
	plex.Metadata.GrandparentTitle = grandParentTitle.String
	plex.Metadata.Index = int(index.Int32)
	plex.Metadata.ParentIndex = int(parent_index.Int32)
	plex.Metadata.Title = title.String

	return &plex, nil
}

func (repo *PlexRepo) Delete(ctx context.Context, req *domain.DeletePlexRequest) error {
	return nil
}
