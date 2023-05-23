package httphelpers

import (
	"context"
	"errors"

	"greenlight/internal/users/models"

	"github.com/gin-gonic/gin"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(ctx *gin.Context, user models.User) {
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx, userContextKey, user))
}

func ContextGetUser(ctx *gin.Context) (models.User, error) {
	user, ok := GetFromContext[models.User](ctx, userContextKey)
	if !ok {
		return models.User{}, errors.New("missing user value in request context")
	}
	return user, nil
}

func GetFromContext[T any](ctx *gin.Context, key any) (T, bool) {
	value := ctx.Request.Context().Value(key)
	if value == nil {
		var zero T
		return zero, false
	}
	ret, ok := value.(T)
	return ret, ok
}
