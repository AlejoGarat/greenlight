package routes

import (
	"expvar"

	"github.com/gin-gonic/gin"
)

type Handler interface{}

func MakeRoutes(engine *gin.RouterGroup) {
	metrics := engine.Group("debug/vars")
	{
		metrics.GET("", gin.WrapH(expvar.Handler()))
	}
}
