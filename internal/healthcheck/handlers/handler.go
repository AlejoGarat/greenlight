package handlers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger  *log.Logger
	Version string
	Env     string
}

func (h *Handler) Healthcheck(c *gin.Context) {
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, h.Env, h.Version)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write([]byte(js))
}
