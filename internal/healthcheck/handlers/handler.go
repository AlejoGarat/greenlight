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
	data := gin.H{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.Env,
			"version":     h.Version,
		},
	}

	err := httphelpers.WriteJSON(c, http.StatusOK, data, nil)
	if err != nil {
		httphelpers.StatusInternalServerErrorResponse(c, err)
	}
}
