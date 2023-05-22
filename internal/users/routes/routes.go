package routes

import (
	"greenlight/internal/users/handlers"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	AddUser() func(c *gin.Context)
	GetUserByEmail() func(c *gin.Context)
	UpdateUser() func(c *gin.Context)
}

func MakeRoutes(engine *gin.RouterGroup, handler *handlers.Handler) {
	users := engine.Group("/users")
	{
		users.POST("", handler.AddUser())
		users.PUT("", handler.UpdateUser())
		users.GET("/:email", handler.GetUserByEmail())
	}
}
