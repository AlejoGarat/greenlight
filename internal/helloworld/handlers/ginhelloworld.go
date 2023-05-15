package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"greenlight/internal/helloworld/logicerrors"
	"greenlight/internal/models"
	"greenlight/pkg/httphelpers"
)

type GinHelloWorldHandler struct {
	logic HelloWorldLogic
}

type HelloWorldLogic interface {
	Greet(ctx context.Context, user *models.User, saveUser bool) (string, error)
	GetUserByName(ctx context.Context, name string) (models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
}

func NewHelloWorldHandler(logic HelloWorldLogic) *GinHelloWorldHandler {
	return &GinHelloWorldHandler{
		logic: logic,
	}
}

// post /helloworld
func (h *GinHelloWorldHandler) Greet() gin.HandlerFunc {
	// you might want to do some processing before returning the handlerFunc,
	// for example if you use a regex, you might want to compile it beforehand
	return func(c *gin.Context) {
		var userInput struct {
			Name string `json:"name"`
		}
		if err := httphelpers.JSONDecode(c, &userInput); err != nil {
			httphelpers.StatusBadRequestResponse(c, err.Error())
			return
		}

		user := models.User{Name: userInput.Name}
		saveUser, err := strconv.ParseBool(c.DefaultQuery("save", "false"))
		if err != nil {
			httphelpers.StatusBadRequestResponse(c, "user not found")
			return
		}

		helloStr, err := h.logic.Greet(c, &user, saveUser)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}

		httphelpers.StatusOKJSONPayloadResponse(c, gin.H{"message": helloStr})
	}
}

// get /helloworld/:name
func (h *GinHelloWorldHandler) GetUserByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		user, err := h.logic.GetUserByName(c, name)
		if err != nil {
			switch {
			case errors.Is(err, logicerrors.ErrUserDoesNotExist):
				httphelpers.StatusBadRequestResponse(c, "user not found")
			default:
				httphelpers.StatusInternalServerErrorResponse(c, err)
			}
			return
		}

		httphelpers.StatusOKJSONPayloadResponse(c, user)
	}
}

// get /helloworld
func (h *GinHelloWorldHandler) ListUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.logic.ListUsers(c)
		if err != nil {
			httphelpers.StatusInternalServerErrorResponse(c, err)
			return
		}

		httphelpers.StatusOKJSONPayloadResponse(c, users)
	}
}
