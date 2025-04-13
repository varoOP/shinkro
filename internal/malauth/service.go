package malauth

import (
	"context"
	"encoding/json"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
	"golang.org/x/oauth2"
)

type Service interface {
	Store(ctx context.Context, ma *domain.MalAuth) error
	Get(ctx context.Context) (*domain.MalAuth, error)
	Delete(ctx context.Context) error
	GetMalClient(ctx context.Context) (*mal.Client, error)
	GetDecrypted(ctx context.Context) (*domain.MalAuth, error)
}

type service struct {
	config *domain.Config
	log    zerolog.Logger
	repo   domain.MalAuthRepo
}

func NewService(config *domain.Config, log zerolog.Logger, repo domain.MalAuthRepo) Service {
	return &service{
		config: config,
		log:    log,
		repo:   repo,
	}
}

func (s *service) Store(ctx context.Context, ma *domain.MalAuth) error {
	et, err := s.config.Encrypt(ma.AccessToken, ma.TokenIV)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to encrypt access token")).Msg("")
		return err
	}

	ecid, err := s.config.Encrypt([]byte(ma.Config.ClientID), ma.TokenIV)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to encrypt client id")).Msg("")
		return err
	}

	ecs, err := s.config.Encrypt([]byte(ma.Config.ClientSecret), ma.TokenIV)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to encrypt client secret")).Msg("")
		return err
	}

	ma.Config.ClientID = string(ecid)
	ma.Config.ClientSecret = string(ecs)
	ma.AccessToken = et
	return s.repo.Store(ctx, ma)
}

func (s *service) Get(ctx context.Context) (*domain.MalAuth, error) {
	return s.repo.Get(ctx)
}

func (s *service) Delete(ctx context.Context) error {
	return s.repo.Delete(ctx)
}

func (s *service) GetDecrypted(ctx context.Context) (*domain.MalAuth, error) {
	ma, err := s.Get(ctx)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to get credentials from database")).Msg("")
		return nil, err
	}

	cid, err := s.config.Decrypt([]byte(ma.Config.ClientID), ma.TokenIV)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to decrypt client id")).Msg("")
		return nil, err
	}

	cs, err := s.config.Decrypt([]byte(ma.Config.ClientSecret), ma.TokenIV)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to decrypt client secret")).Msg("")
		return nil, err
	}

	ma.Config.ClientID = string(cid)
	ma.Config.ClientSecret = string(cs)
	return ma, nil
}

func (s *service) GetMalClient(ctx context.Context) (*mal.Client, error) {
	ma, err := s.GetDecrypted(ctx)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to get credentials from database")).Msg("")
		return nil, err
	}

	token := &oauth2.Token{}
	dt, err := s.config.Decrypt(ma.AccessToken, ma.TokenIV)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to decrypt access token")).Msg("")
		return nil, err
	}

	err = json.Unmarshal(dt, token)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to unmarshal access token")).Msg("")
		return nil, err
	}

	freshToken, err := ma.Config.TokenSource(ctx, token).Token()
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to refresh access token")).Msg("")
		return nil, err
	}

	if freshToken.AccessToken != token.AccessToken {
		token = freshToken
		t, err := json.Marshal(token)
		if err != nil {
			s.log.Err(errors.Wrap(err, "failed to marshal access token")).Msg("")
			return nil, err
		}

		ma.AccessToken = t
		err = s.Store(ctx, ma)
		if err != nil {
			s.log.Err(errors.Wrap(err, "failed to store credentials to database")).Msg("")
			return nil, err
		}
	}

	return mal.NewClient(ma.Config.Client(ctx, freshToken)), nil
}
