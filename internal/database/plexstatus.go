package database

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type PlexStatusRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewPlexStatusRepo(log zerolog.Logger, db *DB) *PlexStatusRepo {
	return &PlexStatusRepo{
		log: log.With().Str("repo", "plex_status").Logger(),
		db:  db,
	}
}

func (repo *PlexStatusRepo) Store(ctx context.Context, ps domain.PlexStatus) error {
	queryBuilder := repo.db.squirrel.
		Insert("plex_status").
		Columns("title", "event", "success", "error_msg", "plex_id").
		Values(ps.Title, ps.Event, ps.Success, ps.ErrorMsg, ps.PlexID).
		Suffix("RETURNING id").RunWith(repo.db.handler)

	var retID int64

	if err := queryBuilder.QueryRowContext(ctx).Scan(&retID); err != nil {
		repo.log.Debug().Err(err).Msg("error executing query")
		return errors.Wrap(err, "error executing query")
	}

	ps.ID = retID
	repo.log.Debug().Msgf("plex.status.store: %+v", ps)
	return nil
}

func (repo *PlexStatusRepo) GetByPlexID(ctx context.Context, plexID int64) (*domain.PlexStatus, error) {
	queryBuilder := repo.db.squirrel.
		Select("id", "title", "event", "success", "error_msg", "plex_id", "time_stamp").
		From("plex_status").
		Where(sq.Eq{"plex_id": plexID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "plex_status.getByPlexID").Msgf("query: '%s', args: '%v'", query, args)

	row := repo.db.handler.QueryRowContext(ctx, query, args...)

	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error querying plex status")
	}

	var ps domain.PlexStatus
	if err := row.Scan(&ps.ID, &ps.Title, &ps.Event, &ps.Success, &ps.ErrorMsg, &ps.PlexID, &ps.TimeStamp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, errors.Wrap(err, "error scanning row")
	}

	return &ps, nil
}

func (repo *PlexStatusRepo) GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]domain.PlexStatus, error) {
	if len(plexIDs) == 0 {
		return []domain.PlexStatus{}, nil
	}

	queryBuilder := repo.db.squirrel.
		Select("id", "title", "event", "success", "error_msg", "plex_id", "time_stamp").
		From("plex_status").
		Where(sq.Eq{"plex_id": plexIDs}).
		OrderBy("time_stamp DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "plex_status.getByPlexIDs").Msgf("query: '%s', args: '%v'", query, args)

	rows, err := repo.db.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	var statuses []domain.PlexStatus
	for rows.Next() {
		var ps domain.PlexStatus
		if err := rows.Scan(&ps.ID, &ps.Title, &ps.Event, &ps.Success, &ps.ErrorMsg, &ps.PlexID, &ps.TimeStamp); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}
		statuses = append(statuses, ps)
	}

	return statuses, nil
}
