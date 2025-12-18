package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
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

func (repo *PlexRepo) Store(ctx context.Context, r *domain.Plex) error {

	guids, err := json.Marshal(r.Metadata.GUID.GUIDS)
	if err != nil {
		return errors.Wrap(err, "failed to marshal GUIDs")
	}

	queryBuilder := repo.db.squirrel.
		Insert("plex_payload").
		Columns("rating", "event", "source", "account_title", "guid_string", "guids", "grand_parent_key", "grand_parent_title", "metadata_index", "library_section_title", "parent_index", "title", "type", "time_stamp").
		Values(r.Rating, r.Event, r.Source, r.Account.Title, r.Metadata.GUID.GUID, string(guids), r.Metadata.GrandparentKey, r.Metadata.GrandparentTitle, r.Metadata.Index, r.Metadata.LibrarySectionTitle, r.Metadata.ParentIndex, r.Metadata.Title, r.Metadata.Type, r.TimeStamp.Format(time.RFC3339)).
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

func (repo *PlexRepo) FindAllWithFilters(ctx context.Context, params domain.PlexPayloadQueryParams) (*domain.FindPlexPayloadsResponse, error) {
	whereQueryBuilder := sq.And{}

	// Apply filters
	if params.Filters.Event != "" {
		whereQueryBuilder = append(whereQueryBuilder, sq.Eq{"p.event": string(params.Filters.Event)})
	}

	if params.Filters.Source != "" {
		whereQueryBuilder = append(whereQueryBuilder, sq.Eq{"p.source": string(params.Filters.Source)})
	}

	// Handle search query - only search by title
	if params.Search != "" {
		search := strings.TrimSpace(params.Search)
		searchPattern := "%" + search + "%"
		whereQueryBuilder = append(whereQueryBuilder, sq.Or{
			sq.Like{"p.title": searchPattern},
			sq.Like{"p.grand_parent_title": searchPattern},
		})
	}

	// Build subquery for pagination
	subQueryBuilder := repo.db.squirrel.
		Select("p.id").
		Distinct().
		From("plex_payload p").
		OrderBy("p.id DESC")

	if params.Limit > 0 {
		subQueryBuilder = subQueryBuilder.Limit(params.Limit)
	} else {
		subQueryBuilder = subQueryBuilder.Limit(20)
	}

	if params.Offset > 0 {
		subQueryBuilder = subQueryBuilder.Offset(params.Offset)
	}

	if len(whereQueryBuilder) > 0 {
		subQueryBuilder = subQueryBuilder.Where(whereQueryBuilder)
	}

	// Handle status filter - need to join with plex_status
	if params.Filters.Status != nil {
		subQueryBuilder = subQueryBuilder.
			InnerJoin("plex_status ps ON p.id = ps.plex_id").
			Where(sq.Eq{"ps.success": *params.Filters.Status})
	}

	subQuery, subArgs, err := subQueryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building subquery")
	}

	// Build count query - reuse whereQueryBuilder
	countQueryBuilder := repo.db.squirrel.
		Select("COUNT(DISTINCT p.id)").
		From("plex_payload p")

	if params.Filters.Status != nil {
		countQueryBuilder = countQueryBuilder.
			InnerJoin("plex_status ps ON p.id = ps.plex_id").
			Where(sq.Eq{"ps.success": *params.Filters.Status})
	}

	if len(whereQueryBuilder) > 0 {
		countQueryBuilder = countQueryBuilder.Where(whereQueryBuilder)
	}

	countQuery, countArgs, err := countQueryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building count query")
	}

	// Build main query
	queryBuilder := repo.db.squirrel.
		Select(
			"p.id", "p.rating", "p.event", "p.source", "p.account_title", "p.guid_string", "p.guids",
			"p.grand_parent_key", "p.grand_parent_title", "p.metadata_index", "p.library_section_title",
			"p.parent_index", "p.title", "p.type", "p.time_stamp",
			"ps.id", "ps.title", "ps.event", "ps.success", "ps.error_type", "ps.error_msg", "ps.plex_id", "ps.time_stamp",
		).
		From("plex_payload p").
		OrderBy("p.id DESC").
		Where("p.id IN ("+subQuery+")", subArgs...).
		LeftJoin("plex_status ps ON p.id = ps.plex_id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	// Execute count query separately
	var totalCount int
	row := repo.db.handler.QueryRowContext(ctx, countQuery, countArgs...)
	if err := row.Scan(&totalCount); err != nil {
		return nil, errors.Wrap(err, "error scanning count")
	}

	allArgs := args

	repo.log.Trace().Str("database", "plex.findAllWithFilters").Msgf("query: '%s', args: '%v'", query, allArgs)

	resp := &domain.FindPlexPayloadsResponse{
		Data:       make([]domain.PlexPayloadListItem, 0),
		TotalCount: totalCount,
	}

	rows, err := repo.db.handler.QueryContext(ctx, query, allArgs...)
	if err != nil {
		return resp, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return resp, errors.Wrap(err, "error rows find plex payloads")
	}

	for rows.Next() {
		var p domain.Plex
		var ps domain.PlexStatus
		var rating sql.NullFloat64
		var index, parent_index sql.NullInt32
		var grandParentKey, grandParentTitle, title, guids, guid_string sql.NullString
		var psID, psPlexID sql.NullInt64
		var psTitle, psEvent, psErrorType, psErrorMsg sql.NullString
		var psSuccess sql.NullBool
		var psTimestamp sql.NullTime

		if err := rows.Scan(
			&p.ID, &rating, &p.Event, &p.Source, &p.Account.Title, &guid_string, &guids,
			&grandParentKey, &grandParentTitle, &index, &p.Metadata.LibrarySectionTitle,
			&parent_index, &title, &p.Metadata.Type, &p.TimeStamp,
			&psID, &psTitle, &psEvent, &psSuccess, &psErrorType, &psErrorMsg, &psPlexID, &psTimestamp,
		); err != nil {
			return resp, errors.Wrap(err, "error scanning row")
		}

		// Unmarshal GUIDs
		if guids.Valid && guids.String != "" {
			err = json.Unmarshal([]byte(guids.String), &p.Metadata.GUID.GUIDS)
			if err != nil {
				return resp, errors.Wrap(err, "error unmarshaling guids")
			}
		}

		p.Metadata.GUID.GUID = guid_string.String
		if rating.Valid {
			p.Rating = float32(rating.Float64)
		}
		p.Metadata.GrandparentKey = grandParentKey.String
		p.Metadata.GrandparentTitle = grandParentTitle.String
		if index.Valid {
			p.Metadata.Index = int(index.Int32)
		}
		if parent_index.Valid {
			p.Metadata.ParentIndex = int(parent_index.Int32)
		}
		p.Metadata.Title = title.String

		// Handle PlexStatus
		var status *domain.PlexStatus
		if psID.Valid {
			ps.ID = psID.Int64
			ps.Title = psTitle.String
			ps.Event = psEvent.String
			ps.Success = psSuccess.Bool
			if psErrorType.Valid {
				ps.ErrorType = domain.PlexErrorType(psErrorType.String)
			}
			ps.ErrorMsg = psErrorMsg.String
			ps.PlexID = psPlexID.Int64
			if psTimestamp.Valid {
				ps.TimeStamp = psTimestamp.Time
			}
			status = &ps
		}

		resp.Data = append(resp.Data, domain.PlexPayloadListItem{
			Plex:   &p,
			Status: status,
		})
	}

	return resp, nil
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

func (repo *PlexRepo) GetRecent(ctx context.Context, limit int) ([]*domain.Plex, error) {
	queryBuilder := repo.db.squirrel.
		Select("id, rating, event, source, account_title, guid_string, guids, grand_parent_key, grand_parent_title, metadata_index, library_section_title, parent_index, title, type, time_stamp").
		From("plex_payload").
		OrderBy("id DESC").
		Limit(uint64(limit))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "plex.getRecent").Msgf("query: '%s', args: '%v'", query, args)

	rows, err := repo.db.handler.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}
	defer rows.Close()

	payloads := make([]*domain.Plex, 0, limit)
	for rows.Next() {
		var p domain.Plex
		var guidsStr string
		if err := rows.Scan(&p.ID, &p.Rating, &p.Event, &p.Source, &p.Account.Title, &p.Metadata.GUID.GUID, &guidsStr, &p.Metadata.GrandparentKey, &p.Metadata.GrandparentTitle, &p.Metadata.Index, &p.Metadata.LibrarySectionTitle, &p.Metadata.ParentIndex, &p.Metadata.Title, &p.Metadata.Type, &p.TimeStamp); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}
		if err := json.Unmarshal([]byte(guidsStr), &p.Metadata.GUID.GUIDS); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling guids")
		}
		payloads = append(payloads, &p)
	}
	return payloads, nil
}
