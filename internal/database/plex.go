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

func (repo *PlexRepo) Store(ctx context.Context, userID int, r *domain.Plex) error {
	r.UserID = userID

	guids, err := json.Marshal(r.Metadata.GUID.GUIDS)
	if err != nil {
		return errors.Wrap(err, "failed to marshal GUIDs")
	}

	queryBuilder := repo.db.squirrel.
		Insert("plex_payload").
		Columns("user_id", "rating", "event", "source", "account_title", "guid_string", "guids", "grand_parent_key", "grand_parent_title", "metadata_index", "library_section_title", "parent_index", "title", "type", "time_stamp").
		Values(r.UserID, r.Rating, r.Event, r.Source, r.Account.Title, r.Metadata.GUID.GUID, string(guids), r.Metadata.GrandparentKey, r.Metadata.GrandparentTitle, r.Metadata.Index, r.Metadata.LibrarySectionTitle, r.Metadata.ParentIndex, r.Metadata.Title, r.Metadata.Type, r.TimeStamp.Format(time.RFC3339)).
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
	//TODO: Implement delete for plex payloads
	return nil
}

func (repo *PlexRepo) CountScrobbleEvents(ctx context.Context) (int, error) {
	queryBuilder := repo.db.squirrel.
		Select("count(*)").
		From("plex_payload").
		Where(sq.Eq{"event": "media.scrobble"})

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

func (repo *PlexRepo) CountRateEvents(ctx context.Context) (int, error) {
	queryBuilder := repo.db.squirrel.
		Select("count(*)").
		From("plex_payload").
		Where(sq.Eq{"event": "media.rate"})

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

func (repo *PlexRepo) GetWithCursor(ctx context.Context, userID int, limit int, cursor *domain.PlexCursor) ([]*domain.Plex, error) {
	queryBuilder := repo.db.squirrel.
		Select("id, user_id, rating, event, source, account_title, guid_string, guids, grand_parent_key, grand_parent_title, metadata_index, library_section_title, parent_index, title, type, time_stamp").
		From("plex_payload").
		Where("user_id = ?", userID)

	if cursor != nil {
		// Page strictly by id to avoid time-zone comparison issues
		queryBuilder = queryBuilder.Where(sq.Lt{"id": cursor.ID})
	}

	// request one extra row to detect presence of next page
	queryBuilder = queryBuilder.OrderBy("id DESC").Limit(uint64(limit + 1))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "plex.getWithCursor").Msgf("query: '%s', args: '%v'", query, args)

	rows, err := repo.db.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	payloads := make([]*domain.Plex, 0)
	for rows.Next() {
		var p domain.Plex
		var guidsStr string
		if err := rows.Scan(&p.ID, &p.UserID, &p.Rating, &p.Event, &p.Source, &p.Account.Title, &p.Metadata.GUID.GUID, &guidsStr, &p.Metadata.GrandparentKey, &p.Metadata.GrandparentTitle, &p.Metadata.Index, &p.Metadata.LibrarySectionTitle, &p.Metadata.ParentIndex, &p.Metadata.Title, &p.Metadata.Type, &p.TimeStamp); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}
		if err := json.Unmarshal([]byte(guidsStr), &p.Metadata.GUID.GUIDS); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling guids")
		}
		payloads = append(payloads, &p)
	}
	return payloads, nil
}

func (repo *PlexRepo) GetWithOffset(ctx context.Context, req *domain.PlexHistoryRequest) ([]*domain.Plex, int, error) {
	// Build base query for counting
	countQueryBuilder := repo.db.squirrel.
		Select("count(*)").
		From("plex_payload").
		Where("user_id = ?", req.UserID)

	// Build base query for data
	queryBuilder := repo.db.squirrel.
		Select("id, user_id, rating, event, source, account_title, guid_string, guids, grand_parent_key, grand_parent_title, metadata_index, library_section_title, parent_index, title, type, time_stamp").
		From("plex_payload").
		Where("user_id = ?", req.UserID)

	// Apply filters
	if req.Search != "" {
		searchFilter := sq.Or{
			sq.Like{"title": "%" + req.Search + "%"},
			sq.Like{"grand_parent_title": "%" + req.Search + "%"},
		}
		countQueryBuilder = countQueryBuilder.Where(searchFilter)
		queryBuilder = queryBuilder.Where(searchFilter)
	}

	if req.Status != "" && req.Status != "all" {
		// This will be handled by joining with plex_status table
		// For now, we'll implement basic filtering
	}

	if req.Event != "" && req.Event != "all" {
		eventFilter := sq.Eq{"event": req.Event}
		countQueryBuilder = countQueryBuilder.Where(eventFilter)
		queryBuilder = queryBuilder.Where(eventFilter)
	}

	if req.FromDate != "" {
		fromDate, err := time.Parse("2006-01-02", req.FromDate)
		if err == nil {
			countQueryBuilder = countQueryBuilder.Where(sq.GtOrEq{"time_stamp": fromDate})
			queryBuilder = queryBuilder.Where(sq.GtOrEq{"time_stamp": fromDate})
		}
	}

	if req.ToDate != "" {
		toDate, err := time.Parse("2006-01-02", req.ToDate)
		if err == nil {
			// Add one day to include the entire day
			toDate = toDate.Add(24 * time.Hour)
			countQueryBuilder = countQueryBuilder.Where(sq.Lt{"time_stamp": toDate})
			queryBuilder = queryBuilder.Where(sq.Lt{"time_stamp": toDate})
		}
	}

	// Get total count
	countQuery, countArgs, err := countQueryBuilder.ToSql()
	if err != nil {
		return nil, 0, errors.Wrap(err, "error building count query")
	}

	repo.log.Trace().Str("database", "plex.getWithOffset.count").Msgf("query: '%s', args: '%v'", countQuery, countArgs)

	var totalCount int
	row := repo.db.handler.QueryRowContext(ctx, countQuery, countArgs...)
	if err := row.Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "error scanning count")
	}

	// Build data query with pagination
	queryBuilder = queryBuilder.OrderBy("time_stamp DESC, id DESC").Limit(uint64(req.Limit)).Offset(uint64(req.Offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "plex.getWithOffset.data").Msgf("query: '%s', args: '%v'", query, args)

	rows, err := repo.db.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	payloads := make([]*domain.Plex, 0)
	for rows.Next() {
		var p domain.Plex
		var guidsStr string
		if err := rows.Scan(&p.ID, &p.UserID, &p.Rating, &p.Event, &p.Source, &p.Account.Title, &p.Metadata.GUID.GUID, &guidsStr, &p.Metadata.GrandparentKey, &p.Metadata.GrandparentTitle, &p.Metadata.Index, &p.Metadata.LibrarySectionTitle, &p.Metadata.ParentIndex, &p.Metadata.Title, &p.Metadata.Type, &p.TimeStamp); err != nil {
			return nil, 0, errors.Wrap(err, "error scanning row")
		}
		if err := json.Unmarshal([]byte(guidsStr), &p.Metadata.GUID.GUIDS); err != nil {
			return nil, 0, errors.Wrap(err, "error unmarshaling guids")
		}
		payloads = append(payloads, &p)
	}

	return payloads, totalCount, nil
}
