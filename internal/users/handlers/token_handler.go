package handlers

import (
	"errors"
	"net/http"
	"time"

	"greenlight/internal/users/models"
	"greenlight/internal/users/serviceerrors"
	"greenlight/pkg/httphelpers"
	"greenlight/pkg/jsonlog"
	"greenlight/pkg/validator"

	"github.com/gin-gonic/gin"
)

type userInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenHandler struct {
	Logger       *jsonlog.Logger
	Version      string
	Env          string
	TokenService TokenService
	UserService  UserService
}

func (h *TokenHandler) CreateAuthToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		var userInput userInput
		err := httphelpers.ReadJSON(c, &userInput)
		if err != nil {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		v := validator.New()
		models.ValidateEmail(v, userInput.Email)
		models.ValidatePasswordPlaintext(v, userInput.Password)

		if !v.Valid() {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		user, err := h.UserService.GetUserByEmail(c, userInput.Email)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrUserNotFound):
				httphelpers.StatusUnauthorizedResponse(c)
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			return
		}

		match, err := user.Password.Matches(userInput.Password)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}
		if !match {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		token, err := h.TokenService.Insert(c, user.ID, 24*time.Hour, models.ScopeAuthentication)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}

		err = httphelpers.WriteJSON(c, http.StatusCreated, gin.H{"authentication_token": token}, nil)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}
	}
}
