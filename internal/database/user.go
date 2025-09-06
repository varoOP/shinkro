package database

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type UserRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewUserRepo(log zerolog.Logger, db *DB) domain.UserRepo {
	return &UserRepo{
		log: log.With().Str("repo", "user").Logger(),
		db:  db,
	}
}

func (r *UserRepo) GetUserCount(ctx context.Context) (int, error) {
	queryBuilder := r.db.squirrel.Select("count(*)").From("users")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "error building query")
	}

	row := r.db.handler.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return 0, errors.Wrap(err, "error executing query")
	}

	result := 0
	if err := row.Scan(&result); err != nil {
		return 0, errors.Wrap(err, "error scanning row")
	}

	return result, nil
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	queryBuilder := r.db.squirrel.
		Select("id", "username", "password", "admin").
		From("users").
		Where(sq.Eq{"username": username})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	row := r.db.handler.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}

	var user domain.User

	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Admin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}

		return nil, errors.Wrap(err, "error scanning row")
	}

	return &user, nil
}

func (r *UserRepo) FindAll(ctx context.Context) ([]*domain.User, error) {
	queryBuilder := r.db.squirrel.
		Select("id", "username", "password", "admin").
		From("users")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	rows, err := r.db.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Admin); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error in rows")
	}

	return users, nil
}

func (r *UserRepo) Store(ctx context.Context, req domain.CreateUserRequest) error {
	queryBuilder := r.db.squirrel.
		Insert("users").
		Columns("username", "password", "admin").
		Values(req.Username, req.Password, req.Admin)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building query")
	}

	_, err = r.db.handler.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing query")
	}

	return err
}

func (r *UserRepo) Update(ctx context.Context, user domain.UpdateUserRequest) error {
	queryBuilder := r.db.squirrel.Update("users")

	if user.UsernameNew != "" {
		queryBuilder = queryBuilder.Set("username", user.UsernameNew)
	}

	if user.PasswordNewHash != "" {
		queryBuilder = queryBuilder.Set("password", user.PasswordNewHash)
	}

	queryBuilder = queryBuilder.Where(sq.Eq{"username": user.UsernameCurrent})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building query")
	}

	_, err = r.db.handler.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing query")
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, username string) error {
	queryBuilder := r.db.squirrel.
		Delete("users").
		Where(sq.Eq{"username": username})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building query")
	}

	// Execute the query.
	_, err = r.db.handler.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing query")
	}

	// Log the deletion.
	r.log.Debug().Msgf("user.delete: successfully deleted user: %s", username)

	return nil
}
