package handlers

import (
	"context"
	"errors"
	"net/http"

	"greenlight/internal/users/models"
	"greenlight/internal/users/serviceerrors"
	"greenlight/pkg/httphelpers"
	"greenlight/pkg/jsonlog"
	"greenlight/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger       *jsonlog.Logger
	Version      string
	Env          string
	UserService  UserService
	TokenService TokenService
}

type UserService interface {
	AddUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	GetForToken(ctx context.Context, tokenScope string, tokenPlaintext string) (models.User, error)
}

type TokenService interface {
	DeleteAllForUser(ctx context.Context, scope string, userID int64) error
}

type createUserInput struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type tokenInput struct {
	TokenPlaintext string `json:"token" db:"name"`
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
			case errors.Is(err, serviceerrors.ErrDuplicatedEmail):
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

func (h *Handler) ActivateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var input tokenInput

		err := httphelpers.ReadJSON(c, &input)
		if err != nil {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		v := validator.New()
		if models.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
			httphelpers.StatusBadRequestJSONPayloadResponse(c, v.Errors)
			return
		}

		user, err := h.UserService.GetForToken(c, models.ScopeActivation, input.TokenPlaintext)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrNoTokenFound):
				v.AddError("token", "invalid or expired activation token")
				httphelpers.StatusBadRequestJSONPayloadResponse(c, v.Errors)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			return
		}

		user.Activated = true

		user, err = h.UserService.UpdateUser(c, user)

		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrEditConflict):
				httphelpers.StatusConflictResponse(c)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			return
		}

		err = h.TokenService.DeleteAllForUser(c, models.ScopeActivation, user.ID)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}

		err = httphelpers.WriteJSON(c, http.StatusOK, gin.H{"user": user}, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
		}
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
