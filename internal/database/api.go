package database

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

func NewAPIRepo(log zerolog.Logger, db *DB) domain.APIRepo {
	return &APIRepo{
		log: log.With().Str("repo", "api").Logger(),
		db:  db,
	}
}

type APIRepo struct {
	log zerolog.Logger
	db  *DB
}

func (r *APIRepo) Store(ctx context.Context, userID int, key *domain.APIKey) error {
	queryBuilder := r.db.squirrel.
		Insert("api_key").
		Columns(
			"user_id",
			"name",
			"key",
			"scopes",
		).
		Values(
			userID,
			key.Name,
			key.Key,
			pq.Array(key.Scopes),
		).
		Suffix("RETURNING created_at").RunWith(r.db.handler)

	var createdAt time.Time

	if err := queryBuilder.QueryRowContext(ctx).Scan(&createdAt); err != nil {
		return errors.Wrap(err, "error executing query")
	}

	key.UserID = userID
	key.CreatedAt = createdAt

	return nil
}

func (r *APIRepo) Delete(ctx context.Context, key string) error {
	userID, err := domain.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	
	queryBuilder := r.db.squirrel.Delete("api_key").Where(sq.Eq{"key": key, "user_id": userID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building query")
	}

	_, err = r.db.handler.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing query")
	}

	r.log.Debug().Msgf("successfully deleted: %v", key)

	return nil
}

func (r *APIRepo) GetAllAPIKeys(ctx context.Context) ([]domain.APIKey, error) {
	userID, err := domain.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	
	queryBuilder := r.db.squirrel.
		Select("user_id", "name", "key", "scopes", "created_at").
		From("api_key").
		Where(sq.Eq{"user_id": userID})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	rows, err := r.db.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}

	defer rows.Close()

	keys := make([]domain.APIKey, 0)
	for rows.Next() {
		var a domain.APIKey

		var name sql.NullString

		if err := rows.Scan(&a.UserID, &name, &a.Key, pq.Array(&a.Scopes), &a.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}

		a.Name = name.String

		keys = append(keys, a)
	}

	return keys, nil
}

func (r *APIRepo) GetKey(ctx context.Context, key string) (*domain.APIKey, error) {
	queryBuilder := r.db.squirrel.
		Select("user_id", "name", "key", "scopes", "created_at").
		From("api_key").
		Where(sq.Eq{"key": key})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	row := r.db.handler.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}

	var apiKey domain.APIKey

	var name sql.NullString

	if err := row.Scan(&apiKey.UserID, &name, &apiKey.Key, pq.Array(&apiKey.Scopes), &apiKey.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}

		return nil, errors.Wrap(err, "error scanning row")
	}

	apiKey.Name = name.String

	return &apiKey, nil
}

func (r *APIRepo) GetUserIDByAPIKey(ctx context.Context, key string) (int, error) {
	queryBuilder := r.db.squirrel.
		Select("user_id").
		From("api_key").
		Where(sq.Eq{"key": key})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "error building query")
	}

	var userID int
	err = r.db.handler.QueryRowContext(ctx, query, args...).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("api key not found")
		}
		return 0, errors.Wrap(err, "error scanning row")
	}

	return userID, nil
}
