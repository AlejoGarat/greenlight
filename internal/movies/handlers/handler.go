package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	models "greenlight/internal/movies/models"
	"greenlight/pkg/httphelpers"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger  *log.Logger
	Version string
	Env     string
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
		fmt.Fprintln(c.Writer, "create a new movie")
	}
}

func (h *Handler) ShowMovie() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := httphelpers.ReadIDParam(c)
		if err != nil {
			http.NotFound(c.Writer, c.Request)
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

		err = httphelpers.WriteJSON(c, http.StatusOK, httphelpers.Envelope{"movie": movie}, nil)
		if err != nil {
			h.Logger.Print(err)
			http.Error(c.Writer, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		}
	}
}
