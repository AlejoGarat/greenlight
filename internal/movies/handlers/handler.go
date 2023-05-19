package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	commonmodels "greenlight/internal/models"
	"greenlight/internal/movies/models"
	"greenlight/internal/serviceerrors"
	"greenlight/pkg/httphelpers"
	"greenlight/pkg/jsonlog"
	"greenlight/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger       *jsonlog.Logger
	Version      string
	Env          string
	MovieService MovieService
}
type createMovieInput struct {
	Title   string               `json:"title"`
	Year    int32                `json:"year"`
	Runtime commonmodels.Runtime `json:"runtime"`
	Genres  []string             `json:"genres"`
}

type updateMovieInput struct {
	Title   *string               `json:"title"`
	Year    *int32                `json:"year"`
	Runtime *commonmodels.Runtime `json:"runtime"`
	Genres  []string              `json:"genres"`
}

type MovieService interface {
	AddMovie(ctx context.Context, movie models.Movie) (models.Movie, error)
	GetMovie(ctx context.Context, id int64) (models.Movie, error)
	GetMovies(ctx context.Context, title string, genres []string, filters commonmodels.Filters) ([]models.Movie, commonmodels.Metadata, error)
	UpdateMovie(ctx context.Context, movie models.Movie) (models.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

func New(logger *jsonlog.Logger, version, env string) *Handler {
	return &Handler{
		Logger:  logger,
		Version: version,
		Env:     env,
	}
}

func (h *Handler) CreateMovie() func(c *gin.Context) {
	return func(c *gin.Context) {
		var input createMovieInput

		err := httphelpers.ReadJSON(c, &input)
		if err != nil {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		movie := models.Movie{
			Title:   input.Title,
			Year:    input.Year,
			Runtime: input.Runtime,
			Genres:  input.Genres,
		}

		v := validator.New()
		valid := fieldsAreValid(c, v, movie)
		if !valid {
			httphelpers.StatusUnprocesableEntities(c, v.Errors)
			return
		}

		ctx := c.Request.Context()
		movie, err = h.MovieService.AddMovie(ctx, movie)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

		err = httphelpers.WriteJSON(c, http.StatusCreated, gin.H{"movie": movie}, headers)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
		}
	}
}

func fieldsAreValid(c *gin.Context, v *validator.Validator, movie models.Movie) bool {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

	return v.Valid()
}

func (h *Handler) ShowMovie() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := httphelpers.ReadIDParam(c)
		if err != nil {
			httphelpers.StatusNotFoundResponse(c)
			return
		}

		ctx := c.Request.Context()
		movie, err := h.MovieService.GetMovie(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrRecordNotFound):
				httphelpers.StatusNotFoundResponse(c)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			return
		}

		err = httphelpers.WriteJSON(c, http.StatusOK, gin.H{"movie": movie}, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
		}
	}
}

func (h *Handler) UpdateMovie() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := httphelpers.ReadIDParam(c)
		if err != nil {
			httphelpers.StatusNotFoundResponse(c)
			return
		}

		ctx := c.Request.Context()
		movie, err := h.MovieService.GetMovie(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrRecordNotFound):
				httphelpers.StatusNotFoundResponse(c)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			return
		}

		var input updateMovieInput
		err = httphelpers.ReadJSON(c, &input)
		if err != nil {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		if input.Title != nil {
			movie.Title = *input.Title
		}

		if input.Year != nil {
			movie.Year = *input.Year
		}

		if input.Runtime != nil {
			movie.Runtime = *input.Runtime
		}

		if input.Genres != nil {
			movie.Genres = input.Genres
		}

		v := validator.New()
		valid := fieldsAreValid(c, v, movie)
		if !valid {
			httphelpers.StatusUnprocesableEntities(c, v.Errors)
			return
		}

		movie, err = h.MovieService.UpdateMovie(ctx, movie)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrEditConflict):
				httphelpers.StatusConflictResponse(c)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}

			return
		}

		err = httphelpers.WriteJSON(c, http.StatusOK, gin.H{"movie": movie}, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
		}
	}
}

func (h *Handler) DeleteMovie() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := httphelpers.ReadIDParam(c)
		if err != nil {
			httphelpers.StatusNotFoundResponse(c)
			return
		}

		ctx := c.Request.Context()
		err = h.MovieService.DeleteMovie(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrRecordNotFound):
				httphelpers.StatusNotFoundResponse(c)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			return
		}

		err = httphelpers.WriteJSON(c, http.StatusOK, gin.H{"message": "movie succesfully deleted"}, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
		}
	}
}

func (h *Handler) ListMovies() func(c *gin.Context) {
	return func(c *gin.Context) {
		var input struct {
			Title  string
			Genres []string
			commonmodels.Filters
		}

		v := validator.New()

		qs := c.Request.URL.Query()

		input.Title = httphelpers.ReadString(qs, "title", "")
		input.Genres = httphelpers.ReadCSV(qs, "genres", []string{})
		input.Filters.Page = httphelpers.ReadInt(qs, "page", 1, v)
		input.Filters.PageSize = httphelpers.ReadInt(qs, "page_size", 20, v)
		input.Filters.Sort = httphelpers.ReadString(qs, "sort", "id")
		input.Filters.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}
		query := struct {
			Query string `json:"query"`
		}{}
		body, _ := io.ReadAll(c.Request.Body)
		json.Unmarshal(body, &query)
		input.Filters.Sort = query.Query

		if commonmodels.ValidateFilters(v, input.Filters); !v.Valid() {
			httphelpers.StatusBadRequestJSONPayloadResponse(c, v.Errors)
			return
		}

		movies, metadata, err := h.MovieService.GetMovies(c.Request.Context(), input.Title, input.Genres, input.Filters)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}

		if len(movies) == 0 {
			movies = []models.Movie{}
		}

		err = httphelpers.WriteJSON(c, http.StatusOK, gin.H{"movies": movies, "metadata": metadata}, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}
	}
}
