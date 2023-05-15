package routes

import (
	"greenlight/internal/healthcheck/handlers"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Healthcheck(c *gin.Context)
}

func MakeRoutes(engine *gin.RouterGroup, handler *handlers.Handler) {
	engine.GET("healthcheck", handler.Healthcheck)
}
