package repo

import (
	"database/sql"
	"errors"

	models "greenlight/internal/movies/models"
	"greenlight/internal/repoerrors"

	"github.com/lib/pq"
)

type movieRepo struct {
	DB *sql.DB
}

func NewMovieRepo(db *sql.DB) *movieRepo {
	return &movieRepo{
		DB: db,
	}
}

func (m movieRepo) Insert(movie *models.Movie) error {
	query := `INSERT INTO movies (title, year, runtime, genres)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m movieRepo) Get(id int64) (*models.Movie, error) {
	if id < 1 {
		return nil, repoerrors.ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, title, year, runtime, genres, version
        FROM movies
        WHERE id = $1`

	var movie models.Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, repoerrors.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m movieRepo) Update(movie *models.Movie) error {
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

	err := m.DB.QueryRow(query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repoerrors.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m movieRepo) Delete(id int64) error {
	if id < 1 {
		return repoerrors.ErrRecordNotFound
	}

	query := `
		DELETE FROM movies
		WHERE id = $1`

	result, err := m.DB.Exec(query, id)
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
