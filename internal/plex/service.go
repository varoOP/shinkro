package plex

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/varoOP/shinkro/internal/notification"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/animeupdate"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
	"github.com/varoOP/shinkro/internal/plexsettings"
	"github.com/varoOP/shinkro/internal/plexstatus"
)

// getUserIDFromContext attempts to get userID from context
// Returns error if not found
func getUserIDFromContext(ctx context.Context) (int, error) {
	if userID, ok := ctx.Value("userID").(int); ok {
		return userID, nil
	}
	return 0, errors.New("userID not found in context")
}

type Service interface {
	Store(ctx context.Context, userID int, plex *domain.Plex) error
	Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error)
	ProcessPlex(ctx context.Context, plex *domain.Plex, agent *domain.PlexSupportedAgents) error
	GetPlexSettings(ctx context.Context) (*domain.PlexSettings, error)
	CountScrobbleEvents(ctx context.Context) (int, error)
	CountRateEvents(ctx context.Context) (int, error)
	GetPlexHistory(ctx context.Context, req *domain.PlexHistoryRequest) (*domain.PlexHistoryResponse, error)
}

type service struct {
	log                 zerolog.Logger
	repo                domain.PlexRepo
	plexettingsService  plexsettings.Service
	animeService        anime.Service
	mapService          mapping.Service
	malauthService      malauth.Service
	animeUpdateService  animeupdate.Service
	notificationService notification.Service
	plexStatusService   plexstatus.Service
}

func NewService(log zerolog.Logger, plexsettingsSvc plexsettings.Service, repo domain.PlexRepo, animeSvc anime.Service, mapSvc mapping.Service, malauthSvc malauth.Service, animeUpdateSvc animeupdate.Service, notificationSvc notification.Service, plexStatusSvc plexstatus.Service) Service {
	return &service{
		log:                 log.With().Str("module", "plex").Logger(),
		repo:                repo,
		plexettingsService:  plexsettingsSvc,
		animeService:        animeSvc,
		mapService:          mapSvc,
		malauthService:      malauthSvc,
		animeUpdateService:  animeUpdateSvc,
		notificationService: notificationSvc,
		plexStatusService:   plexStatusSvc,
	}
}

func (s *service) Get(ctx context.Context, req *domain.GetPlexRequest) (*domain.Plex, error) {
	return s.repo.Get(ctx, req)
}

func (s *service) Store(ctx context.Context, userID int, plex *domain.Plex) error {
	return s.repo.Store(ctx, plex)
}

func (s *service) GetPlexSettings(ctx context.Context) (*domain.PlexSettings, error) {
	// For now, use default user ID 1 when getting plex settings
	// This maintains backward compatibility
	return s.plexettingsService.Get(ctx, 1)
}

func (s *service) ProcessPlex(ctx context.Context, plex *domain.Plex, agent *domain.PlexSupportedAgents) error {
	a, err := s.extractSourceIdForAnime(ctx, plex, agent)
	if err != nil {
		s.plexStatusService.StoreError(ctx, plex, err.Error())
		s.notificationService.Send(domain.NotificationEventError, domain.NotificationPayload{
			Message:      err.Error(),
			Subject:      "Failed to extract anime information",
			AnimeLibrary: plex.Metadata.LibrarySectionTitle,
			PlexEvent:    plex.Event,
			PlexSource:   plex.Source,
			Timestamp:    time.Now(),
		})
		return err
	}

	err = s.animeUpdateService.UpdateAnimeList(ctx, a, plex.Event)
	if err != nil {
		s.plexStatusService.StoreError(ctx, plex, err.Error())
		s.notificationService.Send(domain.NotificationEventError, domain.NotificationPayload{
			Message:      err.Error(),
			Subject:      "Failed to update MyAnimeList",
			AnimeLibrary: a.Plex.Metadata.LibrarySectionTitle,
			PlexEvent:    a.Plex.Event,
			PlexSource:   a.Plex.Source,
			Timestamp:    time.Now(),
		})

		return err
	}

	s.plexStatusService.StoreSuccess(ctx, plex)
	s.notificationService.Send(domain.NotificationEventSuccess, domain.NotificationPayload{
		MediaName:       a.ListDetails.Title,
		MALID:           a.MALId,
		AnimeLibrary:    a.Plex.Metadata.LibrarySectionTitle,
		EpisodesWatched: a.ListStatus.NumEpisodesWatched,
		EpisodesTotal:   a.ListDetails.TotalEpisodeNum,
		TimesRewatched:  a.ListStatus.NumTimesRewatched,
		PictureURL:      a.ListDetails.PictureURL,
		StartDate:       a.ListStatus.StartDate,
		FinishDate:      a.ListStatus.FinishDate,
		AnimeStatus:     string(a.ListStatus.Status),
		Score:           a.ListStatus.Score,
		PlexEvent:       a.Plex.Event,
		PlexSource:      a.Plex.Source,
		Timestamp:       time.Now(),
	})

	return nil
}

func (s *service) extractSourceIdForAnime(ctx context.Context, plex *domain.Plex, agent *domain.PlexSupportedAgents) (*domain.AnimeUpdate, error) {
	source, id, err := s.getSourceIDFromAgent(ctx, plex, agent)
	if err != nil {
		return nil, err
	}

	a := plex.SetAnimeFields(source, id)
	return &a, nil
}

func (s *service) getSourceIDFromAgent(ctx context.Context, p *domain.Plex, agent *domain.PlexSupportedAgents) (domain.PlexSupportedDBs, int, error) {
	switch *agent {
	case domain.HAMA, domain.MALAgent:
		return p.Metadata.GUID.HamaMALAgent(*agent)
	case domain.PlexAgent:
		// For Plex webhooks, we need to get userID from API key context
		// If no userID is found, use default user ID 1
		var userID int = 1
		if ctxUserID, ok := ctx.Value("api_user_id").(int); ok {
			userID = ctxUserID
		}
		return s.plexettingsService.HandlePlexAgent(ctx, userID, p)
	}
	return "", 0, errors.New("unknown agent")
}

func (s *service) CountScrobbleEvents(ctx context.Context) (int, error) {
	return s.repo.CountScrobbleEvents(ctx)
}

func (s *service) CountRateEvents(ctx context.Context) (int, error) {
	return s.repo.CountRateEvents(ctx)
}

func (s *service) GetPlexHistory(ctx context.Context, req *domain.PlexHistoryRequest) (*domain.PlexHistoryResponse, error) {
	// Get Plex payloads based on request type
	var plexPayloads []*domain.Plex
	var totalCount int
	var err error
	var hasMore bool
	var nextCursor string

	if req.Type == "timeline" {
		// Use cursor-based pagination for timeline (with lookahead)
		plexPayloads, err = s.getPlexWithCursor(ctx, req)
		hasMore = len(plexPayloads) > req.Limit
		if hasMore {
			// Determine next cursor from the last item of the current page (index limit-1)
			if req.Limit > 0 && len(plexPayloads) >= req.Limit {
				anchor := plexPayloads[req.Limit-1]
				nextCursor = s.encodeCursor(anchor.TimeStamp, anchor.ID)
			}
			// Trim to page size
			plexPayloads = plexPayloads[:req.Limit]
		} else if len(plexPayloads) > 0 {
			// No more pages, but keep a stable nextCursor empty
			nextCursor = ""
		}
		// totalCount not used for cursor-based
	} else {
		// Use offset-based pagination for table
		plexPayloads, totalCount, err = s.getPlexWithOffset(ctx, req)
	}

	if err != nil {
		return nil, err
	}

	// Extract Plex IDs for batch fetching
	plexIDs := make([]int64, len(plexPayloads))
	for i, p := range plexPayloads {
		plexIDs[i] = p.ID
	}

	// Get PlexStatus for all payloads
	plexStatuses, err := s.plexStatusService.GetByPlexIDs(ctx, plexIDs)
	if err != nil {
		return nil, err
	}

	// Create map for quick lookup
	statusMap := make(map[int64]*domain.PlexStatus)
	for i := range plexStatuses {
		statusMap[plexStatuses[i].PlexID] = &plexStatuses[i]
	}

	// Get AnimeUpdates only for successful ones
	var successfulPlexIDs []int64
	for _, status := range plexStatuses {
		if status.Success {
			successfulPlexIDs = append(successfulPlexIDs, status.PlexID)
		}
	}

	animeUpdates, err := s.getAnimeUpdatesByPlexIDs(ctx, successfulPlexIDs)
	if err != nil {
		return nil, err
	}

	// Create map for quick lookup
	animeUpdateMap := make(map[int64]*domain.AnimeUpdate)
	for _, au := range animeUpdates {
		animeUpdateMap[au.PlexId] = au
	}

	// Combine data
	var items []domain.PlexHistoryItem
	for _, plex := range plexPayloads {
		item := domain.PlexHistoryItem{
			Plex:   plex,
			Status: statusMap[plex.ID],
		}

		if animeUpdate, exists := animeUpdateMap[plex.ID]; exists {
			item.AnimeUpdate = animeUpdate
		}

		items = append(items, item)
	}

	// Build pagination info
	pagination := s.buildPagination(req, len(items), totalCount)

	// For timeline, compute next cursor and hasNext using lookahead
	if req.Type == "timeline" {
		pagination.HasPrev = req.Cursor != ""
		pagination.HasNext = hasMore
		if hasMore && nextCursor != "" {
			pagination.Next = nextCursor
		}
	}

	return &domain.PlexHistoryResponse{
		Data:       items,
		Pagination: pagination,
	}, nil
}

func (s *service) getPlexWithCursor(ctx context.Context, req *domain.PlexHistoryRequest) ([]*domain.Plex, error) {
	var cursor *domain.PlexCursor
	if req.Cursor != "" {
		decoded, err := s.decodeCursor(req.Cursor)
		if err == nil {
			cursor = decoded
		}
	}
	return s.repo.GetWithCursor(ctx, req.Limit, cursor)
}

func (s *service) getPlexWithOffset(ctx context.Context, req *domain.PlexHistoryRequest) ([]*domain.Plex, int, error) {
	payloads, totalCount, err := s.repo.GetWithOffset(ctx, req)
	return payloads, totalCount, err
}

func (s *service) getAnimeUpdatesByPlexIDs(ctx context.Context, plexIDs []int64) ([]*domain.AnimeUpdate, error) {
	if len(plexIDs) == 0 {
		return []*domain.AnimeUpdate{}, nil
	}

	return s.animeUpdateService.GetByPlexIDs(ctx, plexIDs)
}

func (s *service) buildPagination(req *domain.PlexHistoryRequest, itemCount, totalCount int) domain.PlexHistoryPagination {
	pagination := domain.PlexHistoryPagination{}

	if req.Type == "timeline" {
		// Cursor-based pagination
		pagination.HasNext = itemCount == req.Limit
		pagination.HasPrev = req.Cursor != ""
	} else {
		// Offset-based pagination
		pagination.CurrentPage = (req.Offset / req.Limit) + 1
		pagination.TotalPages = (totalCount + req.Limit - 1) / req.Limit
		pagination.TotalItems = totalCount
	}

	return pagination
}

// encodeCursor creates a base64 JSON cursor from timestamp and id
func (s *service) encodeCursor(ts time.Time, id int64) string {
	payload := struct {
		Time string `json:"t"`
		ID   int64  `json:"i"`
	}{Time: ts.UTC().Format(time.RFC3339Nano), ID: id}
	b, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// decodeCursor parses a base64 JSON cursor into PlexCursor
func (s *service) decodeCursor(c string) (*domain.PlexCursor, error) {
	b, err := base64.RawURLEncoding.DecodeString(c)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Time string `json:"t"`
		ID   int64  `json:"i"`
	}
	if err := json.Unmarshal(b, &payload); err != nil {
		return nil, err
	}
	ts, err := time.Parse(time.RFC3339Nano, payload.Time)
	if err != nil {
		return nil, err
	}
	return &domain.PlexCursor{TimeStamp: ts, ID: payload.ID}, nil
}
