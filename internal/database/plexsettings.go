package database

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type PlexSettingsRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewPlexSettingsRepo(log zerolog.Logger, db *DB) domain.PlexSettingsRepo {
	return &PlexSettingsRepo{
		log: log,
		db:  db,
	}
}

func (repo *PlexSettingsRepo) Store(ctx context.Context, ps domain.PlexSettings) (*domain.PlexSettings, error) {

	queryBuilder := repo.db.squirrel.
		Replace("plex_settings").
		Columns("id", "host", "port", "tls", "tls_skip_verify", "token", "username", "anime_libraries", "plex_client_enabled", "client_id").
		Values(0, ps.Host, ps.Port, ps.TLS, ps.TLSSkip, ps.Token, ps.PlexUser, ps.AnimeLibraries, ps.PlexClientEnabled, ps.ClientID).
		RunWith(repo.db.handler)

	_, err := queryBuilder.Exec()
	if err != nil {
		repo.log.Err(err).Msg("error executing query")
		return nil, err
	}

	return &ps, nil
}

func (repo *PlexSettingsRepo) Get(ctx context.Context) (*domain.PlexSettings, error) {
	queryBuilder := repo.db.squirrel.
		Select("ps.host", "ps.port", "ps.tls", "ps.tls_skip_verify", "ps.token", "ps.username", "ps.anime_libraries", "ps.plex_client_enabled", "client_id").
		From("plex_settings ps").
		Where(sq.Eq{"ps.id": 0}).
		RunWith(repo.db.handler)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "plexSettings.get").Msgf("query: '%s', args: '%v'", query, args)
	row := repo.db.handler.QueryRowContext(ctx, query, args...)

	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error rows get malauth")
	}

	var host, token, username, clientID string
	var port int
	var tls, tls_skip_verify, plex_client_enabled bool
	var anime_libraries []string

	if err := row.Scan(&host, &port, &tls, &tls_skip_verify, &token, &username, &anime_libraries, &plex_client_enabled, &clientID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, errors.Wrap(err, "error scanning row")
	}

	ps := domain.NewPlexSettings(host, username, token, clientID, port, anime_libraries, plex_client_enabled, tls, tls_skip_verify)

	return ps, nil
}
