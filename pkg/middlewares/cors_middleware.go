package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EnableCors(trustedOrigins ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodOptions {
			return
		}

		c.Writer.Header().Add("Vary", "Origin")
		c.Writer.Header().Add("Vary", "Access-Control-Request-Method")
		origin := c.Request.Header.Get("Origin")

		for _, trustedOrigin := range trustedOrigins {
			if origin == trustedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				if c.Request.Method == http.MethodOptions {
					c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE")
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
					c.Writer.WriteHeader(http.StatusOK)
					return
				}

				break
			}
		}
		c.AbortWithStatus(http.StatusOK)
	}
}
