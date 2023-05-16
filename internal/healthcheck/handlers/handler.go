package handlers

import (
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

func (h *Handler) Healthcheck(c *gin.Context) {
	data := map[string]string{
		"status":      "available",
		"environment": h.Env,
		"version":     h.Version,
	}

	err := httphelpers.WriteJSON(c, http.StatusOK, httphelpers.Envelope{"data": data}, nil)
	if err != nil {
		h.Logger.Println(err)
		http.Error(c.Writer, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
