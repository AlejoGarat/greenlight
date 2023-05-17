package routes

import (
	"greenlight/internal/movies/handlers"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ShowMovie() func(c *gin.Context)
	CreateMovie() func(c *gin.Context)
	UpdateMovie() func(c *gin.Context)
	DeleteMovie() func(c *gin.Context)
}

func MakeRoutes(engine *gin.RouterGroup, handler *handlers.Handler) {
	engine.GET("movies/:id", handler.ShowMovie())
	engine.POST("movies", handler.CreateMovie())
	engine.PUT("movies/:id", handler.UpdateMovie())
	engine.DELETE("movies/:id", handler.DeleteMovie())
}
