package middlewares

import (
	"github.com/gin-gonic/gin"
)

func EnableCors(trustedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Vary", "Origin")
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			for i := range trustedOrigins {
				if origin == trustedOrigins[i] {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")
		c.Next()
	}
}
