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
	engine.POST("users", handler.AddUser())
	engine.GET("users/:email", handler.GetUserByEmail())
	engine.PUT("users", handler.UpdateUser())
}
