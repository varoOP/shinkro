package malauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
	"golang.org/x/oauth2"
)

type Service interface {
	Store(ctx context.Context, ma *domain.MalAuth) error
	Get(ctx context.Context) (*domain.MalAuth, error)
	GetMalAuthClient(ctx context.Context) (*http.Client, error)
	NewMalAuthClient(ctx context.Context, clientId, clientSecret string) (*domain.MalAuthOpts, error)
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

func (s *service) Get(ctx context.Context) (*domain.MalAuth, error) {
	return s.repo.Get(ctx)
}

func (s *service) GetMalAuthClient(ctx context.Context) (*http.Client, error) {
	ma, err := s.Get(ctx)
	if err != nil {
		return nil, err
	}

	fresh_token, err := ma.Config.TokenSource(ctx, &ma.AccessToken).Token()
	if err != nil {
		return nil, err
	}

	if err == nil && (fresh_token != &ma.AccessToken) {
		ma.AccessToken = *fresh_token
		err = s.Store(ctx, ma)
		if err != nil {
			return nil, err
		}
	}

	return ma.Config.Client(ctx, fresh_token), nil
}

func (s *service) NewMalAuthClient(ctx context.Context, clientId, clientSecret string) (*domain.MalAuthOpts, error) {
	ma := domain.NewMalAuth(clientId, clientSecret)
	verifier, challenge, err := generatePKCE(128)
	if err != nil {
		return nil, err
	}

	codeChallenge := oauth2.SetAuthURLParam("code_challenge", challenge)
	responseType :=  oauth2.SetAuthURLParam("response_type", "code")
	state := randomString(64)

	authCodeUrl := ma.Config.AuthCodeURL(state, codeChallenge, responseType)

	return &domain.MalAuthOpts{
		MalAuth:       ma,
		Verifier:      verifier,
		State:         state,
		AuthCodeUrl: authCodeUrl,
	}, nil
}

func generatePKCE(length int) (verifier, challenge string, err error) {
	if length < 43 || length > 128 {
		return "", "", errors.New("length not supported")
	}

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", errors.Errorf("failed to generate random bytes: %v", err)
	}

	verifier = base64.URLEncoding.EncodeToString(randomBytes)
	verifier = verifier[:length]

	//Waiting for support from MAL side
	// s256 := sha256.New()
	// s256.Write([]byte(verifier))
	// challenge = base64.URLEncoding.EncodeToString(s256.Sum(nil))
	// challenge = base64.RawURLEncoding.EncodeToString(s256.Sum(nil))

	challenge = verifier
	return verifier, challenge, nil
}

func randomString(l int) string {
	random := make([]byte, l)
	_, err := rand.Read(random)
	if err != nil {
		log.Fatalln(err)
	}

	return base64.URLEncoding.EncodeToString(random)[:l]
}
