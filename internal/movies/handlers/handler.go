package handlers

import (
	"fmt"
	"log"
	"net/http"

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

		fmt.Fprintf(c.Writer, "show the details of movie %d\n", id)
	}
}
