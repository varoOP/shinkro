package database

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type AnimeRepo struct {
	log zerolog.Logger
	db  *DB
}

func NewAnimeRepo(log zerolog.Logger, db *DB) domain.AnimeRepo {
	return &AnimeRepo{
		log: log,
		db:  db,
	}
}

func (repo *AnimeRepo) GetByID(ctx context.Context, req *domain.GetAnimeRequest) (*domain.Anime, error) {

	id := "a." + string(req.IDtype) + "_id"

	queryBuilder := repo.db.squirrel.
		Select("a.mal_id", "a.title", "a.en_title", "a.anidb_id", "a.tvdb_id", "a.tmdb_id", "a.type", "a.releaseDate").
		From("anime a").
		Where(sq.Eq{id: req.Id}).
		OrderBy("CASE WHEN a.tvdb_id > 0 THEN 0 ELSE 1 END", "a.mal_id DESC").
		Limit(1)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building query")
	}

	repo.log.Trace().Str("database", "anime.getByTVDBID").Msgf("query: '%s', args: '%v'", query, args)

	row := repo.db.handler.QueryRowContext(ctx, query, args...)

	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "error rows find anime")
	}

	var anime domain.Anime
	var title, enTitle, animeType, releaseDate sql.NullString
	var anidbid, tvdbid, tmdbid sql.NullInt32

	if err := row.Scan(&anime.MALId, &title, &enTitle, &anidbid, &tvdbid, &tmdbid, &animeType, &releaseDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, errors.Wrap(err, "error scanning row")
	}

	anime.MainTitle = title.String
	anime.EnglishTitle = enTitle.String
	anime.AnimeType = animeType.String
	anime.ReleaseDate = releaseDate.String
	anime.AniDBId = int(anidbid.Int32)
	anime.TVDBId = int(tvdbid.Int32)
	anime.TMDBId = int(tmdbid.Int32)

	return &anime, nil
}

func (repo *AnimeRepo) StoreMultiple(anime []*domain.Anime) error {
	tx, err := repo.db.handler.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	for _, a := range anime {
		queryBuilder := repo.db.squirrel.
			Replace("anime").
			Columns("mal_id", "title", "en_title", "anidb_id", "tvdb_id", "tmdb_id", "type", "releaseDate").
			Values(a.MALId, a.MainTitle, a.EnglishTitle, a.AniDBId, a.TVDBId, a.TMDBId, a.AnimeType, a.ReleaseDate).
			RunWith(tx)

		_, err := queryBuilder.Exec()
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
