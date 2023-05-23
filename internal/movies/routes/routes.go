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
	ListMovies() func(c *gin.Context)
}

func MakeRoutes(engine *gin.RouterGroup, handler *handlers.Handler) {
	movies := engine.Group("movies")
	{
		movies.GET("", handler.ListMovies())
		movies.GET("/:id", handler.ShowMovie())
		movies.POST("", handler.CreateMovie())
		movies.PATCH("/:id", handler.UpdateMovie())
		movies.DELETE("/:id", handler.DeleteMovie())
	}
}
