package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	commonmodels "greenlight/internal/models"
	models "greenlight/internal/movies/models"
	"greenlight/pkg/httphelpers"
	"greenlight/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger  *log.Logger
	Version string
	Env     string
}
type movieInput struct {
	Title   string               `json:"title"`
	Year    int32                `json:"year"`
	Runtime commonmodels.Runtime `json:"runtime"`
	Genres  []string             `json:"genres"`
}

type MovieService interface {
	AddMovie(movie *models.Movie) error
	GetMovie(id int64) (*models.Movie, error)
	UpdateMovie(movie *models.Movie) error
	DeleteMovie(id int64) error
}

func New(logger *log.Logger, version, env string) *Handler {
	return &Handler{
		Logger:  logger,
		Version: version,
		Env:     env,
	}
}

func (h *Handler) CreateMovie() func(c *gin.Context) {
	return func(c *gin.Context) {
		var movie movieInput

		err := httphelpers.ReadJSON(c, &movie)
		if err != nil {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		v := validator.New()
		validateFields(c, v, movie)

		fmt.Fprintf(c.Writer, "%+v\n", movie)
	}
}

func validateFields(c *gin.Context, v *validator.Validator, input movieInput) {
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	if !v.Valid() {
		httphelpers.StatusUnprocesableEntities(c, v.Errors)
		return
	}
}

func (h *Handler) ShowMovie() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := httphelpers.ReadIDParam(c)
		if err != nil {
			httphelpers.StatusNotFoundResponse(c)
			return
		}

		movie := models.Movie{
			ID:        id,
			CreatedAt: time.Now(),
			Title:     "Casablanca",
			Runtime:   102,
			Genres:    []string{"drama", "romance", "war"},
			Version:   1,
		}

		data := gin.H{
			"movie": movie,
		}

		err = httphelpers.WriteJSON(c, http.StatusOK, data, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
		}
	}
}
