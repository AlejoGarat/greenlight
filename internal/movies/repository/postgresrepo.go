package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	commonmodels "greenlight/internal/models"
	models "greenlight/internal/movies/models"
	"greenlight/internal/repoerrors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type movieRepo struct {
	DB *sqlx.DB
}

func NewMovieRepo(db *sqlx.DB) *movieRepo {
	return &movieRepo{
		DB: db,
	}
}

func (m movieRepo) Insert(ctx context.Context, movie models.Movie) (models.Movie, error) {
	query := `INSERT INTO movies (title, year, runtime, genres)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := m.DB.GetContext(ctx, &movie, query, args...)
	if err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}

func (m movieRepo) Get(ctx context.Context, id int64) (models.Movie, error) {
	if id < 1 {
		return models.Movie{}, repoerrors.ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, title, year, runtime, genres, version
        FROM movies
        WHERE id = $1`

	var movie models.Movie

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := m.DB.GetContext(ctx, &movie, query, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNoRows):
			return models.Movie{}, repoerrors.ErrRecordNotFound
		default:
			return models.Movie{}, err
		}
	}

	return movie, nil
}

func (m movieRepo) GetAll(ctx context.Context, title string, genres []string, filters commonmodels.Filters) ([]models.Movie, error) {
	column, err := filters.SortColumn()
	if err != nil {
		return []models.Movie{}, err
	}

	query := fmt.Sprintf(`
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
		AND (genres @> $2 OR $2 = '{}')     
		ORDER BY %s %s, id ASC`, column, filters.SortDirection())

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var movies []models.Movie

	err = m.DB.SelectContext(ctx, &movies, query, title, pq.StringArray(genres))
	if err != nil {
		return movies, err
	}

	return movies, nil
}

func (m movieRepo) Update(ctx context.Context, movie models.Movie) (models.Movie, error) {
	query := `
        UPDATE movies 
        SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := m.DB.GetContext(ctx, &movie, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Movie{}, repoerrors.ErrEditConflict
		default:
			return models.Movie{}, err
		}
	}

	return movie, nil
}

func (m movieRepo) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return repoerrors.ErrRecordNotFound
	}

	query := `
		DELETE FROM movies
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repoerrors.ErrRecordNotFound
	}

	return nil
}
