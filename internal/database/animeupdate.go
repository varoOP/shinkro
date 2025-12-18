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
		Columns("mal_id", "source_db", "source_id", "episode_num", "season_num", "time_stamp", "list_details", "list_status", "plex_id", "status", "error_type", "error_message").
		Values(r.MALId, r.SourceDB, r.SourceId, r.EpisodeNum, r.SeasonNum, r.Timestamp, string(listDetails), string(listStatus), r.PlexId, r.Status, r.ErrorType, r.ErrorMessage).
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

func (repo *AnimeUpdateRepo) GetRecentUnique(ctx context.Context, limit int) ([]*domain.AnimeUpdate, error) {
	latest := repo.db.squirrel.
		Select("mal_id, MAX(time_stamp) AS max_ts").
		From("anime_update").
		GroupBy("mal_id")

	queryBuilder := repo.db.squirrel.
		Select("au.id, au.mal_id, au.source_db, au.source_id, au.episode_num, au.season_num, au.time_stamp, au.list_details, au.list_status, au.plex_id, au.status, au.error_type, au.error_message").
		FromSelect(latest, "latest").
		Join("anime_update au ON latest.mal_id = au.mal_id AND latest.max_ts = au.time_stamp").
		OrderBy("au.time_stamp DESC").
		Limit(uint64(limit))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	rows, err := repo.db.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	updates := make([]*domain.AnimeUpdate, 0)
	for rows.Next() {
		var au domain.AnimeUpdate
		var listDetailsBytes, listStatusBytes []byte
		var status, errorType, errorMessage sql.NullString
		if err := rows.Scan(&au.ID, &au.MALId, &au.SourceDB, &au.SourceId, &au.EpisodeNum, &au.SeasonNum, &au.Timestamp, &listDetailsBytes, &listStatusBytes, &au.PlexId, &status, &errorType, &errorMessage); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}
		if err := json.Unmarshal(listDetailsBytes, &au.ListDetails); err != nil {
			return nil, errors.Wrap(err, "error unmarshalling list_details")
		}
		if err := json.Unmarshal(listStatusBytes, &au.ListStatus); err != nil {
			return nil, errors.Wrap(err, "error unmarshalling list_status")
		}
		if status.Valid {
			au.Status = domain.AnimeUpdateStatusType(status.String)
		}
		if errorType.Valid {
			au.ErrorType = domain.AnimeUpdateErrorType(errorType.String)
		}
		if errorMessage.Valid {
			au.ErrorMessage = errorMessage.String
		}
		updates = append(updates, &au)
	}
	return updates, nil
}

func (repo *AnimeUpdateRepo) GetByPlexID(ctx context.Context, plexID int64) (*domain.AnimeUpdate, error) {
	queryBuilder := repo.db.squirrel.
		Select("id, mal_id, source_db, source_id, episode_num, season_num, time_stamp, list_details, list_status, plex_id, status, error_type, error_message").
		From("anime_update").
		Where("plex_id = ?", plexID).
		OrderBy("time_stamp DESC").
		Limit(1)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	row := repo.db.handler.QueryRowContext(ctx, query, args...)
	var au domain.AnimeUpdate
	var listDetailsBytes, listStatusBytes []byte
	var status, errorType, errorMessage sql.NullString
	if err := row.Scan(&au.ID, &au.MALId, &au.SourceDB, &au.SourceId, &au.EpisodeNum, &au.SeasonNum, &au.Timestamp, &listDetailsBytes, &listStatusBytes, &au.PlexId, &status, &errorType, &errorMessage); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No update for this plex_id
		}
		return nil, errors.Wrap(err, "error scanning row")
	}
	if err := json.Unmarshal(listDetailsBytes, &au.ListDetails); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling list_details")
	}
	if err := json.Unmarshal(listStatusBytes, &au.ListStatus); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling list_status")
	}
	if status.Valid {
		au.Status = domain.AnimeUpdateStatusType(status.String)
	}
	if errorType.Valid {
		au.ErrorType = domain.AnimeUpdateErrorType(errorType.String)
	}
	if errorMessage.Valid {
		au.ErrorMessage = errorMessage.String
	}
	return &au, nil
}

func (repo *AnimeUpdateRepo) GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]*domain.AnimeUpdate, error) {
	if len(plexIDs) == 0 {
		return []*domain.AnimeUpdate{}, nil
	}

	queryBuilder := repo.db.squirrel.
		Select("id, mal_id, source_db, source_id, episode_num, season_num, time_stamp, list_details, list_status, plex_id, status, error_type, error_message").
		From("anime_update").
		Where(sq.Eq{"plex_id": plexIDs}).
		OrderBy("time_stamp DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "animeupdate.getByPlexIDs").Msgf("query: '%s', args: '%v'", query, args)

	rows, err := repo.db.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	var updates []*domain.AnimeUpdate
	for rows.Next() {
		var au domain.AnimeUpdate
		var listDetailsBytes, listStatusBytes []byte
		var status, errorType, errorMessage sql.NullString
		if err := rows.Scan(&au.ID, &au.MALId, &au.SourceDB, &au.SourceId, &au.EpisodeNum, &au.SeasonNum, &au.Timestamp, &listDetailsBytes, &listStatusBytes, &au.PlexId, &status, &errorType, &errorMessage); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}
		if err := json.Unmarshal(listDetailsBytes, &au.ListDetails); err != nil {
			return nil, errors.Wrap(err, "error unmarshalling list_details")
		}
		if err := json.Unmarshal(listStatusBytes, &au.ListStatus); err != nil {
			return nil, errors.Wrap(err, "error unmarshalling list_status")
		}
		if status.Valid {
			au.Status = domain.AnimeUpdateStatusType(status.String)
		}
		if errorType.Valid {
			au.ErrorType = domain.AnimeUpdateErrorType(errorType.String)
		}
		if errorMessage.Valid {
			au.ErrorMessage = errorMessage.String
		}
		updates = append(updates, &au)
	}

	return updates, nil
}
