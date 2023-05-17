package repo

import (
	models "greenlight/internal/movies/models"

	"github.com/jmoiron/sqlx"
)

type movieRepo struct {
	DB *sqlx.DB
}

func NewMovieRepo(db sqlx.DB) *movieRepo {
	return &movieRepo{
		DB: &db,
	}
}

func (m movieRepo) Insert(movie *models.Movie) error {
	return nil
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
