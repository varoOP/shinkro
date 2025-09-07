package database

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type MappingRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewMappingRepo(log zerolog.Logger, db *DB) domain.MappingRepo {
	return &MappingRepo{
		log: log,
		db:  db,
	}
}

func (repo *MappingRepo) Store(ctx context.Context, userID int, m *domain.MapSettings) error {

	queryBuilder := repo.db.squirrel.
		Replace("mapping_settings").
		Columns("user_id", "tvdb_enabled", "tmdb_enabled", "tvdb_path", "tmdb_path").
		Values(userID, m.TVDBEnabled, m.TMDBEnabled, m.CustomMapTVDBPath, m.CustomMapTMDBPath).
		RunWith(repo.db.handler)

	_, err := queryBuilder.Exec()
	if err != nil {
		repo.log.Err(err).Msg("error executing query")
		return err
	}

	return nil
}

func (repo *MappingRepo) Get(ctx context.Context) (*domain.MapSettings, error) {
	userID, err := domain.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	
	queryBuilder := repo.db.squirrel.
		Select("m.user_id", "m.tvdb_enabled", "m.tmdb_enabled", "m.tvdb_path", "m.tmdb_path").
		From("mapping_settings m").
		Where(sq.Eq{"m.user_id": userID}).
		RunWith(repo.db.handler)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "mapping.get").Msgf("query: '%s', args: '%v'", query, args)
	row := repo.db.handler.QueryRowContext(ctx, query, args...)

	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error rows get mapping")
	}

	var tvdbPath, tmdbPath string
	var tvdb, tmdb bool
	var dbUserID int

	if err := row.Scan(&dbUserID, &tvdb, &tmdb, &tvdbPath, &tmdbPath); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, errors.Wrap(err, "error scanning row")
	}

	m := domain.NewMapSettings(tvdb, tmdb, tvdbPath, tmdbPath)
	m.UserID = dbUserID

	return m, nil
}
