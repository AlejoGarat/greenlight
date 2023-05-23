package middlewares

import (
	"context"
	"errors"
	"strings"

	"greenlight/internal/movies/repoerrors"
	"greenlight/internal/users/models"
	"greenlight/pkg/httphelpers"
	"greenlight/pkg/validator"

	"github.com/gin-gonic/gin"
)

type UserRepo interface {
	GetForToken(ctx context.Context, tokenScope string, tokenPlaintext string) (models.User, error)
}

func Authenticate(userRepo UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Authorization")

		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			httphelpers.ContextSetUser(c, models.AnonymousUser)
			c.Next()
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			httphelpers.StatusUnauthorizedResponse(c)
			c.Abort()
			return
		}

		token := headerParts[1]

		v := validator.New()
		if models.ValidateTokenPlaintext(v, token); !v.Valid() {
			httphelpers.StatusUnauthorizedResponse(c)
			c.Abort()
			return
		}

		user, err := userRepo.GetForToken(c, models.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, repoerrors.ErrRecordNotFound):
				httphelpers.StatusUnauthorizedResponse(c)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			c.Abort()
			return
		}

		httphelpers.ContextSetUser(c, user)
	}
}
