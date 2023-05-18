package service

import (
	"context"

	commonmodels "greenlight/internal/models"
	models "greenlight/internal/movies/models"
)

type movieService struct {
	repo MovieRepo
}
type MovieRepo interface {
	Insert(ctx context.Context, movie models.Movie) (models.Movie, error)
	Get(ctx context.Context, id int64) (models.Movie, error)
	GetAll(ctx context.Context, title string, genres []string, filters commonmodels.Filters) ([]models.Movie, error)
	Update(ctx context.Context, movie models.Movie) (models.Movie, error)
	Delete(ctx context.Context, id int64) error
}

func NewMovieService(repo MovieRepo) *movieService {
	return &movieService{
		repo: repo,
	}
}

func (m movieService) AddMovie(ctx context.Context, movie models.Movie) (models.Movie, error) {
	movie, err := m.repo.Insert(ctx, movie)
	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

func (m movieService) GetMovie(ctx context.Context, id int64) (models.Movie, error) {
	movie, err := m.repo.Get(ctx, id)
	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

func (m movieService) GetMovies(ctx context.Context, title string, genres []string, filters commonmodels.Filters) ([]models.Movie, error) {
	movies, err := m.repo.GetAll(ctx, title, genres, filters)
	if err != nil {
		return []models.Movie{}, err
	}

	return movies, nil
}

func (m movieService) UpdateMovie(ctx context.Context, movie models.Movie) (models.Movie, error) {
	movie, err := m.repo.Update(ctx, movie)
	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

func (m movieService) DeleteMovie(ctx context.Context, id int64) error {
	err := m.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
