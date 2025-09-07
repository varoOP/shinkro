package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	List(ctx context.Context) ([]domain.APIKey, error)
	Store(ctx context.Context, key *domain.APIKey) error
	Delete(ctx context.Context, key string) error
	ValidateAPIKey(ctx context.Context, token string) bool
	GetUserIDByAPIKey(ctx context.Context, token string) (int, error)
}

type service struct {
	log  zerolog.Logger
	repo domain.APIRepo

	keyCache map[string]domain.APIKey
}

func NewService(log zerolog.Logger, repo domain.APIRepo) Service {
	return &service{
		log:      log.With().Str("module", "api").Logger(),
		repo:     repo,
		keyCache: map[string]domain.APIKey{},
	}
}

func (s *service) List(ctx context.Context) ([]domain.APIKey, error) {
	userID, err := domain.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(s.keyCache) > 0 {
		keys := make([]domain.APIKey, 0, len(s.keyCache))

		for _, key := range s.keyCache {
			// Filter by userID when returning from cache
			if key.UserID == userID {
				keys = append(keys, key)
			}
		}

		return keys, nil
	}

	return s.repo.GetAllAPIKeys(ctx, userID)
}

func (s *service) Store(ctx context.Context, apiKey *domain.APIKey) error {
	userID, err := domain.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	
	apiKey.Key = GenerateSecureToken(16)

	if err := s.repo.Store(ctx, userID, apiKey); err != nil {
		return err
	}

	if len(s.keyCache) > 0 {
		// set new apiKey
		s.keyCache[apiKey.Key] = *apiKey
	}

	return nil
}

func (s *service) Delete(ctx context.Context, key string) error {
	userID, err := domain.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	
	_, err = s.repo.GetKey(ctx, key)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, userID, key)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not delete api key: %s", key))
	}

	// remove key from cache
	delete(s.keyCache, key)

	return nil
}

func (s *service) ValidateAPIKey(ctx context.Context, key string) bool {
	if _, ok := s.keyCache[key]; ok {
		s.log.Trace().Msgf("api service key cache hit: %s", key)
		return true
	}

	apiKey, err := s.repo.GetKey(ctx, key)
	if err != nil {
		s.log.Trace().Msgf("api service key cache invalid key: %s", key)
		return false
	}

	s.log.Trace().Msgf("api service key cache miss: %s", key)

	s.keyCache[key] = *apiKey

	return true
}

func (s *service) GetUserIDByAPIKey(ctx context.Context, token string) (int, error) {
	// First check cache
	if apiKey, ok := s.keyCache[token]; ok {
		return apiKey.UserID, nil
	}

	// If not in cache, get from database
	return s.repo.GetUserIDByAPIKey(ctx, token)
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
