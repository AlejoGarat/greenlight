package httphelpers

import (
	"context"
	"errors"

	"greenlight/internal/users/models"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(ctx context.Context, user models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func ContextGetUser(ctx context.Context) (models.User, error) {
	user, ok := ctx.Value(userContextKey).(models.User)
	if !ok {
		return models.User{}, errors.New("missing user value in request context")
	}
	return user, nil
}
