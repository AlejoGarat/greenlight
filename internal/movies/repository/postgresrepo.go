package repo

import (
	"database/sql"

	models "greenlight/internal/movies/models"

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
	return nil, nil
}

func (m movieRepo) Update(movie *models.Movie) error {
	return nil
}

func (m movieRepo) Delete(id int64) error {
	return nil
}
