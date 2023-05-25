package service

import (
	"context"
	"errors"

	commonmodels "greenlight/internal/models"
	"greenlight/internal/movies/models"
	"greenlight/internal/movies/repoerrors"
	"greenlight/internal/movies/serviceerrors"
)

type movieService struct {
	repo MovieRepo
}
type MovieRepo interface {
	Insert(ctx context.Context, movie models.Movie) (models.Movie, error)
	Get(ctx context.Context, id int64) (models.Movie, error)
	GetAll(ctx context.Context, title string, genres []string, filters commonmodels.Filters,
	) ([]models.Movie, commonmodels.Metadata, error)
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
		switch {
		case errors.Is(err, repoerrors.ErrMovieTitleRequired):
			return models.Movie{}, serviceerrors.ErrMovieTitleRequired
		case errors.Is(err, repoerrors.ErrMovieYearRequired):
			return models.Movie{}, serviceerrors.ErrMovieYearRequired
		default:
			return models.Movie{}, err
		}
	}

	return movie, nil
}

func (m movieService) GetMovie(ctx context.Context, id int64) (models.Movie, error) {
	movie, err := m.repo.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrMovieNoFound):
			return models.Movie{}, serviceerrors.ErrNoMovieFound
		default:
			return models.Movie{}, err
		}
	}

	return movie, nil
}

func (m movieService) GetMovies(ctx context.Context, title string, genres []string, filters commonmodels.Filters,
) ([]models.Movie, commonmodels.Metadata, error) {
	movies, metadata, err := m.repo.GetAll(ctx, title, genres, filters)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrMovieNoFound):
			return []models.Movie{}, commonmodels.Metadata{}, serviceerrors.ErrNoMovieFound
		default:
			return []models.Movie{}, commonmodels.Metadata{}, err
		}
	}

	return movies, metadata, nil
}

func (m movieService) UpdateMovie(ctx context.Context, movie models.Movie) (models.Movie, error) {
	movie, err := m.repo.Update(ctx, movie)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrMovieNoFound):
			return models.Movie{}, serviceerrors.ErrNoMovieFound
		case errors.Is(err, repoerrors.ErrMovieTitleRequired):
			return models.Movie{}, serviceerrors.ErrMovieTitleRequired
		case errors.Is(err, repoerrors.ErrMovieYearRequired):
			return models.Movie{}, serviceerrors.ErrMovieYearRequired
		default:
			return models.Movie{}, err
		}
	}

	return movie, nil
}

func (m movieService) DeleteMovie(ctx context.Context, id int64) error {
	err := m.repo.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrMovieNoFound):
			return serviceerrors.ErrNoMovieFound
		default:
			return err
		}
	}

	return nil
}
