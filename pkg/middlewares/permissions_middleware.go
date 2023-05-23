package middlewares

import (
	"greenlight/pkg/httphelpers"

	"github.com/gin-gonic/gin"
)

func RequireActivatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := httphelpers.ContextGetUser(c)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			c.Abort()
			return
		}

		if user.IsAnonymous() {
			httphelpers.StatusForbiddenJSONPayloadResponse(c, map[string]any{"error": "auth needed"})
			c.Abort()

			return
		}

		if !user.Activated {
			httphelpers.StatusForbiddenJSONPayloadResponse(c, map[string]any{"error": "activation required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
