package handlers

import (
	"context"

	"greenlight/internal/users/models"
	"greenlight/pkg/jsonlog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger      *jsonlog.Logger
	Version     string
	Env         string
	UserService UserService
}

type UserService interface {
	AddUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
}

func New(logger *jsonlog.Logger, version, env string) *Handler {
	return &Handler{
		Logger:  logger,
		Version: version,
		Env:     env,
	}
}

func (h *Handler) AddUser() func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

func (h *Handler) GetUserByEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

func (h *Handler) UpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}
