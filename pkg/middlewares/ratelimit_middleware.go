package middlewares

import (
	"greenlight/pkg/httphelpers"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(2, 4)

	return func(c *gin.Context) {
		defer func() {
			if !limiter.Allow() {
				httphelpers.RateLimitExceededResponse(c)
				return
			}
		}()
		c.Next()
	}
}
