package database

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type AnimeUpdateRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewAnimeUpdateRepo(log zerolog.Logger, db *DB) domain.AnimeUpdateRepo {
	return &AnimeUpdateRepo{
		log: log.With().Str("repo", "animeupdate").Logger(),
		db:  db,
	}
}

func (repo *AnimeUpdateRepo) Store(ctx context.Context, r *domain.AnimeUpdate) error {
	listDetails, err := json.Marshal(r.ListDetails)
	if err != nil {
		return errors.Wrap(err, "failed to marshal listDetails")
	}

	listStatus, err := json.Marshal(r.ListStatus)
	if err != nil {
		return errors.Wrap(err, "failed to marshal listStatus")
	}

	queryBuilder := repo.db.squirrel.
		Insert("anime_update").
		Columns("mal_id", "source_db", "source_id", "episode_num", "season_num", "time_stamp", "list_details", "list_status", "plex_id").
		Values(r.MALId, r.SourceDB, r.SourceId, r.EpisodeNum, r.SeasonNum, r.Timestamp, string(listDetails), string(listStatus), r.PlexId).
		Suffix("RETURNING id").RunWith(repo.db.handler)

	var retID int64

	if err := queryBuilder.QueryRowContext(ctx).Scan(&retID); err != nil {
		repo.log.Debug().Err(err).Msg("error executing query")
		return errors.Wrap(err, "error executing query")
	}

	r.ID = retID
	repo.log.Debug().Msgf("animeUpdate.store: %+v", r)
	return nil
}

func (repo *AnimeUpdateRepo) GetByID(ctx context.Context, req *domain.GetAnimeUpdateRequest) (*domain.AnimeUpdate, error) {
	return nil, nil
}

func (repo *AnimeUpdateRepo) Count(ctx context.Context) (int, error) {
	queryBuilder := repo.db.squirrel.
		Select("count(*)").
		From("anime_update")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "error building query")
	}

	row := repo.db.handler.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return 0, errors.Wrap(err, "error executing query")
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, errors.Wrap(err, "error scanning row")
	}

	return count, nil
}
