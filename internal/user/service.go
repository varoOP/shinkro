package user

import (
	"context"
	"errors"
	"github.com/rs/zerolog"

	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	GetUserCount(ctx context.Context) (int, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindAll(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, req domain.CreateUserRequest) error
	CreateUserAdmin(ctx context.Context, req domain.CreateUserRequest) error
	Update(ctx context.Context, req domain.UpdateUserRequest) error
	Delete(ctx context.Context, username string) error
}

type service struct {
	repo domain.UserRepo
	log  zerolog.Logger
}

func NewService(repo domain.UserRepo, log zerolog.Logger) Service {
	return &service{
		repo: repo,
		log:  log.With().Str("module", "user").Logger(),
	}
}

func (s *service) GetUserCount(ctx context.Context) (int, error) {
	return s.repo.GetUserCount(ctx)
}

func (s *service) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) FindAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Remove passwords from response
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}

func (s *service) CreateUser(ctx context.Context, req domain.CreateUserRequest) error {
	userCount, err := s.repo.GetUserCount(ctx)
	if err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("only 1 user account is supported at the moment")
	}

	// First user is always admin
	req.Admin = true

	return s.repo.Store(ctx, req)
}

func (s *service) CreateUserAdmin(ctx context.Context, req domain.CreateUserRequest) error {
	// Admin can create users without the limit
	return s.repo.Store(ctx, req)
}

func (s *service) Update(ctx context.Context, req domain.UpdateUserRequest) error {
	return s.repo.Update(ctx, req)
}

func (s *service) Delete(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, username)
}
