package middlewares

import (
	"context"

	"greenlight/internal/permissions/models"
	"greenlight/pkg/httphelpers"

	"github.com/gin-gonic/gin"
)

type PermissionsRepo interface {
	GetAllForUser(ctx context.Context, userID int64) (models.Permissions, error)
}

func RequireAuthenticatedUser(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := httphelpers.ContextGetUser(c)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			c.Abort()
			return
		}

		if user.IsAnonymous() {
			httphelpers.StatusUnauthorizedResponse(c)
			c.Abort()
			return
		}

		next(c)
	}
}

func RequireActivatedUser(next gin.HandlerFunc) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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

	return RequireAuthenticatedUser(fn)
}

func RequirePermission(permissionsRepo PermissionsRepo, code string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		user, err := httphelpers.ContextGetUser(c)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			c.Abort()
			return
		}

		permissions, err := permissionsRepo.GetAllForUser(c, user.ID)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			c.Abort()
			return
		}

		if !permissions.Include(code) {
			httphelpers.StatusForbiddenResponse(c)
			c.Abort()
			return
		}

		c.Next()
	}

	return RequireActivatedUser(fn)
}
