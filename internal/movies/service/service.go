package service

import (
	models "greenlight/internal/movies/models"
)

type movieService struct {
	repo MovieRepo
}
type MovieRepo interface {
	Insert(movie *models.Movie) error
	Get(id int64) (*models.Movie, error)
	Update(movie *models.Movie) error
	Delete(id int64) error
}

func NewMovieService(repo MovieRepo) *movieService {
	return &movieService{
		repo: repo,
	}
}

func (m movieService) AddMovie(movie *models.Movie) error {
	err := m.repo.Insert(movie)
	if err != nil {
		return err
	}

	return nil
}

func (m movieService) GetMovie(id int64) (*models.Movie, error) {
	movie, err := m.repo.Get(id)
	if err != nil {
		return movie, err
	}

	return nil, nil
}

func (m movieService) UpdateMovie(movie *models.Movie) error {
	err := m.repo.Update(movie)
	if err != nil {
		return err
	}

	return nil
}

func (m movieService) DeleteMovie(id int64) error {
	err := m.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
