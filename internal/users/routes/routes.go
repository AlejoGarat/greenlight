package routes

import (
	"greenlight/internal/users/handlers"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	AddUser() func(c *gin.Context)
	GetUserByEmail() func(c *gin.Context)
	UpdateUser() func(c *gin.Context)
	ActivateUser() func(c *gin.Context)
}
type THandler interface {
	CreateAuthToken() func(c *gin.Context)
}

func MakeRoutes(engine *gin.RouterGroup, handler *handlers.Handler, thandler *handlers.TokenHandler) {
	users := engine.Group("/users")
	{
		users.POST("", handler.AddUser())
		users.PUT("", handler.UpdateUser())
		users.PUT("/activated", handler.ActivateUser())
		users.GET("/:email", handler.GetUserByEmail())
	}

	tokens := engine.Group("/tokens")
	{
		tokens.POST("/authentication", thandler.CreateAuthToken())
	}
}
