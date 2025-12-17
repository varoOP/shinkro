package database

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type AnimeUpdateStatusRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewAnimeUpdateStatusRepo(log zerolog.Logger, db *DB) domain.AnimeUpdateStatusRepo {
	return &AnimeUpdateStatusRepo{
		log: log.With().Str("repo", "animeupdatestatus").Logger(),
		db:  db,
	}
}

func (repo *AnimeUpdateStatusRepo) Store(ctx context.Context, status *domain.AnimeUpdateStatus) error {
	queryBuilder := repo.db.squirrel.
		Insert("anime_update_status").
		Columns("plex_id", "mal_id", "status", "error_type", "error_message", "anime_title", "source_db", "source_id", "season_num", "episode_num", "time_stamp").
		Values(
			status.PlexID,
			status.MALID,
			status.Status,
			status.ErrorType,
			status.ErrorMessage,
			status.AnimeTitle,
			status.SourceDB,
			status.SourceID,
			status.SeasonNum,
			status.EpisodeNum,
			status.Timestamp,
		).
		Suffix("RETURNING id").RunWith(repo.db.handler)

	var retID int64
	if err := queryBuilder.QueryRowContext(ctx).Scan(&retID); err != nil {
		repo.log.Debug().Err(err).Msg("error executing query")
		return errors.Wrap(err, "error executing query")
	}

	status.ID = retID
	return nil
}

func (repo *AnimeUpdateStatusRepo) GetByPlexID(ctx context.Context, plexID int64) (*domain.AnimeUpdateStatus, error) {
	queryBuilder := repo.db.squirrel.
		Select("id", "plex_id", "mal_id", "status", "error_type", "error_message", "anime_title", "source_db", "source_id", "season_num", "episode_num", "time_stamp").
		From("anime_update_status").
		Where(sq.Eq{"plex_id": plexID}).
		OrderBy("time_stamp DESC").
		Limit(1).
		RunWith(repo.db.handler)

	var status domain.AnimeUpdateStatus
	var errorType, errorMessage, animeTitle, sourceDB sql.NullString
	var malID, sourceID, seasonNum, episodeNum sql.NullInt64

	err := queryBuilder.QueryRowContext(ctx).Scan(
		&status.ID,
		&status.PlexID,
		&malID,
		&status.Status,
		&errorType,
		&errorMessage,
		&animeTitle,
		&sourceDB,
		&sourceID,
		&seasonNum,
		&episodeNum,
		&status.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		repo.log.Debug().Err(err).Msg("error executing query")
		return nil, errors.Wrap(err, "error executing query")
	}

	if malID.Valid {
		status.MALID = int(malID.Int64)
	}
	if errorType.Valid {
		status.ErrorType = domain.AnimeUpdateErrorType(errorType.String)
	}
	if errorMessage.Valid {
		status.ErrorMessage = errorMessage.String
	}
	if animeTitle.Valid {
		status.AnimeTitle = animeTitle.String
	}
	if sourceDB.Valid {
		status.SourceDB = domain.PlexSupportedDBs(sourceDB.String)
	}
	if sourceID.Valid {
		status.SourceID = int(sourceID.Int64)
	}
	if seasonNum.Valid {
		status.SeasonNum = int(seasonNum.Int64)
	}
	if episodeNum.Valid {
		status.EpisodeNum = int(episodeNum.Int64)
	}

	return &status, nil
}

func (repo *AnimeUpdateStatusRepo) GetByPlexIDs(ctx context.Context, plexIDs []int64) ([]domain.AnimeUpdateStatus, error) {
	if len(plexIDs) == 0 {
		return []domain.AnimeUpdateStatus{}, nil
	}

	queryBuilder := repo.db.squirrel.
		Select("id", "plex_id", "mal_id", "status", "error_type", "error_message", "anime_title", "source_db", "source_id", "season_num", "episode_num", "time_stamp").
		From("anime_update_status").
		Where(sq.Eq{"plex_id": plexIDs}).
		OrderBy("time_stamp DESC").
		RunWith(repo.db.handler)

	rows, err := queryBuilder.QueryContext(ctx)
	if err != nil {
		repo.log.Debug().Err(err).Msg("error executing query")
		return nil, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	var statuses []domain.AnimeUpdateStatus
	for rows.Next() {
		var status domain.AnimeUpdateStatus
		var errorType, errorMessage, animeTitle, sourceDB sql.NullString
		var malID, sourceID, seasonNum, episodeNum sql.NullInt64

		if err := rows.Scan(
			&status.ID,
			&status.PlexID,
			&malID,
			&status.Status,
			&errorType,
			&errorMessage,
			&animeTitle,
			&sourceDB,
			&sourceID,
			&seasonNum,
			&episodeNum,
			&status.Timestamp,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}

		if malID.Valid {
			status.MALID = int(malID.Int64)
		}
		if errorType.Valid {
			status.ErrorType = domain.AnimeUpdateErrorType(errorType.String)
		}
		if errorMessage.Valid {
			status.ErrorMessage = errorMessage.String
		}
		if animeTitle.Valid {
			status.AnimeTitle = animeTitle.String
		}
		if sourceDB.Valid {
			status.SourceDB = domain.PlexSupportedDBs(sourceDB.String)
		}
		if sourceID.Valid {
			status.SourceID = int(sourceID.Int64)
		}
		if seasonNum.Valid {
			status.SeasonNum = int(seasonNum.Int64)
		}
		if episodeNum.Valid {
			status.EpisodeNum = int(episodeNum.Int64)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

