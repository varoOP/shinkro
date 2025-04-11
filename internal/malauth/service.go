package malauth

import (
	"context"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	Store(ctx context.Context, ma *domain.MalAuth) error
	StoreMalAuthOpts(ctx context.Context, mo *domain.MalAuthOpts) error
	get(ctx context.Context) (*domain.MalAuth, error)
	GetMalAuthOpts(ctx context.Context) (*domain.MalAuthOpts, error)
	GetMalClient(ctx context.Context) (*mal.Client, error)
}

type service struct {
	log  zerolog.Logger
	repo domain.MalAuthRepo
}

func NewService(log zerolog.Logger, repo domain.MalAuthRepo) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) Store(ctx context.Context, ma *domain.MalAuth) error {
	return s.repo.Store(ctx, ma)
}

func (s *service) StoreMalAuthOpts(ctx context.Context, mo *domain.MalAuthOpts) error {
	return s.repo.StoreMalAuthOpts(ctx, mo)
}

func (s *service) get(ctx context.Context) (*domain.MalAuth, error) {
	return s.repo.Get(ctx)
}

func (s *service) GetMalAuthOpts(ctx context.Context) (*domain.MalAuthOpts, error) {
	return s.repo.GetMalAuthOpts(ctx)
}

func (s *service) GetMalClient(ctx context.Context) (*mal.Client, error) {
	ma, err := s.get(ctx)
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to get credentials from database")).Msg("")
		return nil, err
	}

	fresh_token, err := ma.Config.TokenSource(ctx, &ma.AccessToken).Token()
	if err != nil {
		s.log.Err(errors.Wrap(err, "failed to refresh access token")).Msg("")
		return nil, err
	}

	if fresh_token.AccessToken != ma.AccessToken.AccessToken {
		ma.AccessToken = *fresh_token
		err = s.Store(ctx, ma)
		if err != nil {
			s.log.Err(errors.Wrap(err, "failed to store credentials to database")).Msg("")
			return nil, err
		}
	}

	return mal.NewClient(ma.Config.Client(ctx, fresh_token)), nil
}
