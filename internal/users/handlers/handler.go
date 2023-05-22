package handlers

import (
	"context"
	"errors"
	"net/http"

	"greenlight/internal/serviceerrors"
	"greenlight/internal/users/models"
	"greenlight/pkg/httphelpers"
	"greenlight/pkg/jsonlog"
	"greenlight/pkg/validator"

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

type createUserInput struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
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
		var userInput createUserInput

		err := httphelpers.ReadJSON(c, &userInput)
		if err != nil {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		user := models.User{
			Name:      userInput.Name,
			Email:     userInput.Email,
			Activated: false,
		}

		err = user.Password.Set(userInput.Password)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}

		v := validator.New()
		if !fieldsAreValid(c, v, user) {
			httphelpers.StatusUnprocesableEntities(c, v.Errors)
			return
		}

		user, err = h.UserService.AddUser(c, user)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrDuplicateEmail):
				v.AddError("email", "a user with this email address already exists")
				httphelpers.StatusUnprocesableEntities(c, v.Errors)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			return
		}

		err = httphelpers.WriteJSON(c, http.StatusAccepted, gin.H{"user": user}, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
		}
	}
}

func fieldsAreValid(c *gin.Context, v *validator.Validator, user models.User) bool {
	models.ValidateEmail(v, user.Email)
	models.ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	models.ValidateUser(v, &user)
	return v.Valid()
}

func (h *Handler) GetUserByEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}

func (h *Handler) UpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
	}
}
