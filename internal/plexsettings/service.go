package plexsettings

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/dcarbone/zadapters/zstdlog"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/pkg/plex"
)

type Service interface {
	Store(ctx context.Context, ps domain.PlexSettings) (*domain.PlexSettings, error)
	Get(ctx context.Context) (*domain.PlexSettings, error)
	Delete(ctx context.Context) error
	GetClient(ctx context.Context) (*plex.Client, error)
	GetEncryptionKey() ([]byte, error)
	HandlePlexAgent(ctx context.Context, p *domain.Plex) (domain.PlexSupportedDBs, int, error)
}

type service struct {
	config *domain.Config
	log    zerolog.Logger
	repo   domain.PlexSettingsRepo
}

func NewService(config *domain.Config, log zerolog.Logger, repo domain.PlexSettingsRepo) Service {
	return &service{
		config: config,
		log:    log,
		repo:   repo,
	}
}

func (s *service) Store(ctx context.Context, ps domain.PlexSettings) (*domain.PlexSettings, error) {
	return s.repo.Store(ctx, ps)
}

func (s *service) Get(ctx context.Context) (*domain.PlexSettings, error) {
	return s.repo.Get(ctx)
}

func (s *service) Delete(ctx context.Context) error {
	return s.repo.Delete(ctx)
}

func (s *service) GetClient(ctx context.Context) (*plex.Client, error) {
	ps, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	scheme := "http"
	if ps.TLS {
		scheme = "https"
	}

	if ps.Host == "" {
		ps.Host = "localhost"
	}

	if ps.Port == 0 {
		ps.Port = 32400
	}

	if len(ps.Token) == 0 || len(ps.TokenIV) == 0 {
		return nil, errors.New("token or tokenIV is empty")
	}

	key, err := s.GetEncryptionKey()
	if err != nil {
		return nil, err
	}

	token, err := decryptToken(ps.Token, key, ps.TokenIV)
	if err != nil {
		return nil, err
	}

	c := plex.NewClient(plex.Config{
		Url:           fmt.Sprintf("%s://%s:%d", scheme, ps.Host, ps.Port),
		Token:         token,
		ClientID:      ps.ClientID,
		TLSSkipVerify: ps.TLSSkip,
		Log:           zstdlog.NewStdLoggerWithLevel(s.log.With().Str("client", "plex").Logger(), zerolog.TraceLevel),
	})

	return c, nil
}

func (s *service) HandlePlexAgent(ctx context.Context, p *domain.Plex) (domain.PlexSupportedDBs, int, error) {
	if p.Metadata.Type == domain.PlexEpisode {
		pc, err := s.GetClient(ctx)
		if err != nil {
			return "", 0, err
		}

		guid, err := pc.GetShowID(ctx, p.Metadata.GrandparentKey)
		if err != nil {
			return "", 0, err
		}

		id := domain.GUID{
			GUIDS: guid.GUIDS,
			GUID:  guid.GUID,
		}

		return id.PlexAgent(p.Metadata.Type)
	}
	return "", 0, nil
}

func (s *service) GetEncryptionKey() ([]byte, error) {
	key, err := LoadKeyFromHex(s.config.EncryptionKey)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func decryptToken(ciphertext []byte, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func LoadKeyFromHex(hexKey string) ([]byte, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, errors.New("invalid hex encryption key")
	}
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}
	return key, nil
}
