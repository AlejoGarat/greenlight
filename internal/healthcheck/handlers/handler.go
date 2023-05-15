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
	fmt.Fprintln(c.Writer, "status: available")
	fmt.Fprintf(c.Writer, "environment: %s\n", h.Env)
	fmt.Fprintf(c.Writer, "version: %s\n", h.Version)
}
