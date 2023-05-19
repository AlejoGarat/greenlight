package middlewares

import (
	"fmt"

	"greenlight/pkg/httphelpers"

	"github.com/gin-gonic/gin"
)

func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if errPanic := recover(); errPanic != nil {
				err, ok := errPanic.(error)
				if !ok {
					err = fmt.Errorf("%v", errPanic)
				}
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
		}()
		c.Next()
	}
}
