package middlewares

import (
	"greenlight/pkg/httphelpers"

	"github.com/gin-gonic/gin"
)

func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httphelpers.StatusInternalServerErrorResponse(c, err.(error))
			}
		}()
		c.Next()
	}
}
