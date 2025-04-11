package database

import (
	"context"
	"database/sql"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type MalAuthRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewMalAuthRepo(log zerolog.Logger, db *DB) domain.MalAuthRepo {
	return &MalAuthRepo{
		log: log,
		db:  db,
	}
}

func (repo *MalAuthRepo) StoreMalAuthOpts(ctx context.Context, mo *domain.MalAuthOpts) error {
	queryBuilder := repo.db.squirrel.
		Replace("malauth").
		Columns("id", "client_id", "client_secret").
		Values(1, mo.ClientID, mo.ClientSecret).
		RunWith(repo.db.handler)

	_, err := queryBuilder.Exec()
	if err != nil {
		repo.log.Err(err).Msg("error executing query")
		return err
	}

	return nil
}

func (repo *MalAuthRepo) GetMalAuthOpts(ctx context.Context) (*domain.MalAuthOpts, error) {
	queryBuilder := repo.db.squirrel.
		Select("ma.client_id", "ma.client_secret").
		From("malauth ma").
		Where(sq.Eq{"ma.id": 1}).
		RunWith(repo.db.handler)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "malauth.get").Msgf("query: '%s', args: '%v'", query, args)
	row := repo.db.handler.QueryRowContext(ctx, query, args...)

	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error rows get malauth")
	}

	var clientId, clientSecret string

	if err := row.Scan(&clientId, &clientSecret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, errors.Wrap(err, "error scanning row")
	}

	return domain.NewMalAuthOpts(clientId, clientSecret, "", "", ""), nil
}

func (repo *MalAuthRepo) Store(ctx context.Context, ma *domain.MalAuth) error {

	accessToken, err := json.Marshal(ma.AccessToken)
	if err != nil {
		repo.log.Err(err).Msg("unable to marshal access token.")
		return err
	}

	queryBuilder := repo.db.squirrel.
		Replace("malauth").
		Columns("id", "client_id", "client_secret", "access_token").
		Values(ma.Id, ma.Config.ClientID, ma.Config.ClientSecret, string(accessToken)).
		RunWith(repo.db.handler)

	_, err = queryBuilder.Exec()
	if err != nil {
		repo.log.Err(err).Msg("error executing query")
		return err
	}

	return nil
}

func (repo *MalAuthRepo) Get(ctx context.Context) (*domain.MalAuth, error) {
	queryBuilder := repo.db.squirrel.
		Select("ma.client_id", "ma.client_secret", "ma.access_token").
		From("malauth ma").
		Where(sq.Eq{"ma.id": 1}).
		RunWith(repo.db.handler)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "malauth.get").Msgf("query: '%s', args: '%v'", query, args)
	row := repo.db.handler.QueryRowContext(ctx, query, args...)

	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error rows get malauth")
	}

	var accessToken, clientId, clientSecret string

	if err := row.Scan(&clientId, &clientSecret, &accessToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, errors.Wrap(err, "error scanning row")
	}

	ma := domain.NewMalAuth(clientId, clientSecret)
	err = json.Unmarshal([]byte(accessToken), &ma.AccessToken)
	if err != nil {
		return nil, err
	}

	return ma, nil
}
