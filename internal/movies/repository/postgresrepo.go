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
	return nil
}

func (m movieRepo) Delete(id int64) error {
	return nil
}
