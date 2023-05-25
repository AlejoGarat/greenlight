package handlers

import (
	"context"
	"net/http"

	"greenlight/pkg/httphelpers"
	"greenlight/pkg/jsonlog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger   *jsonlog.Logger
	Version  string
	Env      string
	DBStatus func(context.Context) error
}

type systemInfo struct {
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type databaseInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type health struct {
	Status     string       `json:"status"`
	SystemInfo systemInfo   `json:"system_info"`
	Database   databaseInfo `json:"database"`
}

func (h *Handler) Healthcheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		data := health{
			Status: "available",
			SystemInfo: systemInfo{
				Environment: h.Env,
				Version:     h.Version,
			},
			Database: databaseInfo{
				Status: "available",
			},
		}

		ok := true
		if err := h.DBStatus(c.Request.Context()); err != nil {
			data.Database.Status = "unavailable"
			data.Database.Message = err.Error()
			ok = false
		}

		status := http.StatusOK
		if !ok {
			status = http.StatusBadRequest
		}
		err := httphelpers.WriteJSON(c, status, data, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
		}
	}
}
