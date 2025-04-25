package anime

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type Service interface {
	GetByID(ctx context.Context, req *domain.GetAnimeRequest) (*domain.Anime, error)
	StoreMultiple(anime []*domain.Anime) error
	GetAnime() ([]*domain.Anime, error)
	UpdateAnime() error
}

type service struct {
	log  zerolog.Logger
	repo domain.AnimeRepo
}

func NewService(log zerolog.Logger, repo domain.AnimeRepo) Service {
	return &service{
		log:  log.With().Str("module", "anime").Logger(),
		repo: repo,
	}
}

func (s *service) GetByID(ctx context.Context, req *domain.GetAnimeRequest) (*domain.Anime, error) {
	return s.repo.GetByID(ctx, req)
}

func (s *service) StoreMultiple(anime []*domain.Anime) error {
	return s.repo.StoreMultiple(anime)
}

func (s *service) GetAnime() ([]*domain.Anime, error) {
	resp, err := http.Get(string(domain.ShinkroDB))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from shinkrodb")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	a := []*domain.Anime{}
	err = json.Unmarshal(body, &a)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal json")
	}

	return a, nil
}

func (s *service) UpdateAnime() error {
	s.log.Info().Msg("Started anime update in database.")
	anime, err := s.GetAnime()
	if err != nil {
		return err
	}

	err = s.StoreMultiple(anime)
	if err != nil {
		return err
	}

	s.log.Info().Msg("Anime update complete.")

	return nil
}
