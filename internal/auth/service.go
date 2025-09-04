package auth

import (
	"context"

	"github.com/alexedwards/argon2id"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/user"
)

type Service interface {
	GetUserCount(ctx context.Context) (int, error)
	Login(ctx context.Context, username, password string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindAll(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, req domain.CreateUserRequest) error
	CreateUserAdmin(ctx context.Context, req domain.CreateUserRequest) error
	UpdateUser(ctx context.Context, req domain.UpdateUserRequest) error
	Delete(ctx context.Context, username string) error
	ResetPassword(ctx context.Context, username, newPassword string) error
	CreateHash(password string) (hash string, err error)
	ComparePasswordAndHash(password string, hash string) (match bool, err error)
}

type service struct {
	log     zerolog.Logger
	userSvc user.Service
}

func NewService(log zerolog.Logger, userSvc user.Service) Service {
	return &service{
		log:     log.With().Str("module", "auth").Logger(),
		userSvc: userSvc,
	}
}

func (s *service) GetUserCount(ctx context.Context) (int, error) {
	return s.userSvc.GetUserCount(ctx)
}

func (s *service) Login(ctx context.Context, username, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("empty credentials supplied")
	}

	// find user
	u, err := s.userSvc.FindByUsername(ctx, username)
	if err != nil {
		s.log.Error().Err(err).Msgf("could not find user by username: %v", username)
		return nil, errors.Wrapf(err, "invalid login: %s", username)
	}

	if u == nil {
		return nil, errors.Errorf("invalid login: %s", username)
	}

	// compare password from request and the saved password
	match, err := s.ComparePasswordAndHash(password, u.Password)
	if err != nil {
		return nil, errors.New("error checking credentials")
	}

	if !match {
		s.log.Error().Msg("bad credentials")
		return nil, errors.Errorf("invalid login: %s", username)
	}

	return u, nil
}

func (s *service) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userSvc.FindByUsername(ctx, username)
}

func (s *service) CreateUser(ctx context.Context, req domain.CreateUserRequest) error {
	if req.Username == "" {
		return errors.New("validation error: empty username supplied")
	} else if req.Password == "" {
		return errors.New("validation error: empty password supplied")
	}

	userCount, err := s.userSvc.GetUserCount(ctx)
	if err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("only 1 user account is supported at the moment")
	}

	hashed, err := s.CreateHash(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	req.Password = hashed

	if err := s.userSvc.CreateUser(ctx, req); err != nil {
		s.log.Error().Err(err).Msgf("could not create user: %s", req.Username)
		return errors.New("failed to create new user")
	}

	return nil
}

func (s *service) UpdateUser(ctx context.Context, req domain.UpdateUserRequest) error {
	if req.PasswordCurrent == "" {
		return errors.New("validation error: empty current password supplied")
	}

	if req.PasswordNew != "" && req.PasswordCurrent != "" {
		if req.PasswordNew == req.PasswordCurrent {
			return errors.New("validation error: new password must be different")
		}
	}

	// find user
	u, err := s.userSvc.FindByUsername(ctx, req.UsernameCurrent)
	if err != nil {
		s.log.Trace().Err(err).Msgf("invalid login %v", req.UsernameCurrent)
		return errors.Wrapf(err, "invalid login: %s", req.UsernameCurrent)
	}

	if u == nil {
		return errors.Errorf("invalid login: %s", req.UsernameCurrent)
	}

	// compare password from request and the saved password
	match, err := s.ComparePasswordAndHash(req.PasswordCurrent, u.Password)
	if err != nil {
		return errors.New("error checking credentials")
	}

	if !match {
		s.log.Debug().Msgf("bad credentials: %q | %q", req.UsernameCurrent, req.PasswordCurrent)
		return errors.Errorf("invalid login: %s", req.UsernameCurrent)
	}

	if req.PasswordNew != "" {
		hashed, err := s.CreateHash(req.PasswordNew)
		if err != nil {
			return errors.New("failed to hash password")
		}

		req.PasswordNewHash = hashed
	}

	if err := s.userSvc.Update(ctx, req); err != nil {
		s.log.Error().Err(err).Msgf("could not change password for user: %s", req.UsernameCurrent)
		return errors.New("failed to change password")
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, username, newPassword string) error {
	if username == "" {
		return errors.New("validation error: empty username supplied")
	}
	if newPassword == "" {
		return errors.New("validation error: empty new password supplied")
	}

	// find user
	u, err := s.userSvc.FindByUsername(ctx, username)
	if err != nil {
		s.log.Trace().Err(err).Msgf("could not find user: %v", username)
		return errors.Wrapf(err, "user not found: %s", username)
	}

	if u == nil {
		return errors.Errorf("user not found: %s", username)
	}

	// hash the new password
	hashed, err := s.CreateHash(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// create update request with only the new password
	updateReq := domain.UpdateUserRequest{
		UsernameCurrent: username,
		PasswordNew:     newPassword,
		PasswordNewHash: hashed,
	}

	if err := s.userSvc.Update(ctx, updateReq); err != nil {
		s.log.Error().Err(err).Msgf("could not reset password for user: %s", username)
		return errors.New("failed to reset password")
	}

	return nil
}

func (s *service) ComparePasswordAndHash(password string, hash string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func (s *service) CreateHash(password string) (hash string, err error) {
	if password == "" {
		return "", errors.New("must supply non empty password to CreateHash")
	}

	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (s *service) FindAll(ctx context.Context) ([]*domain.User, error) {
	return s.userSvc.FindAll(ctx)
}

func (s *service) CreateUserAdmin(ctx context.Context, req domain.CreateUserRequest) error {
	// Hash password before storing
	hash, err := s.CreateHash(req.Password)
	if err != nil {
		return errors.Wrap(err, "could not create password hash")
	}
	req.Password = hash

	return s.userSvc.CreateUserAdmin(ctx, req)
}

func (s *service) Delete(ctx context.Context, username string) error {
	return s.userSvc.Delete(ctx, username)
}
