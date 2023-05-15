package httphelpers

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContentType string

var (
	ErrInvalidPayloadType = errors.New("invalid payload type")
	ErrUnknownContentType = errors.New("unknown content type")
)

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeXML  ContentType = "application/xml"
	ContentTypeHTML ContentType = "text/html"
)

// StatusOKResponse sets an empty 200 response
func StatusOKResponse(c *gin.Context) {
	c.Status(http.StatusOK)
}

// StatusCreatedResponse sets a 201 response and loads a JSON payload containing `{"id":id}`
func StatusCreatedResponse[T int | int64 | string](c *gin.Context, id T) {
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// StatusNoContentResponse sets an empy 204 response
func StatusNoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// StatusBadRequestResponse sets a 400 response and loads a JSON payload containing
// `{"error":"msg"}“
func StatusBadRequestResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

// StatusUnauthorizedResponse sets a 401 response and loads a JSON payload containing
// `{"error":"unauthorized"}“
func StatusUnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

// StatusForbiddenResponse sets a 403 response and loads a JSON payload containing
// `{"error":"forbidden"}“
func StatusForbiddenResponse(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
}

// StatusNotFoundResponse sets a 404 response and loads a JSON payload containing
// `{"error":"not found"}“
func StatusNotFoundResponse(c *gin.Context) {
	CustomStatusJSONPayloadResponse(c, http.StatusNotFound,
		map[string]string{"error": "not found"})
}

// StatusConflictResponse sets a 409 response and loads a JSON payload containing
// `{"error":"the resource you are trying to edit has been modified by another user, please try again"}“
func StatusConflictResponse(c *gin.Context) {
	c.JSON(http.StatusConflict, gin.H{"error": "the resource you are trying to edit has been modified by another user, please try again"})
}

// StatusUnprocesableEntities sets a 422 response and loads a payload containing the errors
func StatusUnprocesableEntities(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
}

// StatusInternalServerErrorResponse sets an empty 500 response and loads errors into context,
// in order to be accessible to middlewares
func StatusInternalServerErrorResponse(c *gin.Context, err error) {
	c.Error(err)
	c.Status(http.StatusInternalServerError)
}

// StatusOKJSONPayloadResponse is a shorthand for CustomStatusJSONPayloadResponse with status 200
//
// If you do not wish to handle the error, and are ok with a 400 response on error, feel free to use c.JSON(http.StatusOK, payload)
func StatusOKJSONPayloadResponse(c *gin.Context, payload any) error {
	return CustomStatusJSONPayloadResponse(c, http.StatusOK, payload)
}

// StatusCreatedJSONPayloadResponse is a shorthand for CustomStatusJSONPayloadResponse with status 201
//
// If you do not wish to handle the error, and are ok with a 400 response on error, feel free to use c.JSON(http.StatusCreated, payload)
func StatusCreatedJSONPayload(c *gin.Context, payload any) error {
	return CustomStatusJSONPayloadResponse(c, http.StatusCreated, payload)
}

// StatusBadRequestJSONPayloadResponse is a shorthand for CustomStatusJSONPayloadResponse with status 400
//
// If you do not wish to handle the error, and are ok with a 400 response on error, feel free to use c.JSON(http.StatusBadRequest, payload)
func StatusBadRequestJSONPayloadResponse(c *gin.Context, payload any) error {
	return CustomStatusJSONPayloadResponse(c, http.StatusBadRequest, payload)
}

// StatusUnauthorizedJSONPayloadResponse is a shorthand for CustomStatusJSONPayloadResponse with status 401
//
// If you do not wish to handle the error, and are ok with a 400 response on error, feel free to use c.JSON(http.StatusUnauthorized, payload)
func StatusUnauthorizedJSONPayloadResponse(c *gin.Context, payload any) error {
	return CustomStatusJSONPayloadResponse(c, http.StatusUnauthorized, payload)
}

// StatusJSONPayloadResponse is a shorthand for CustomStatusJSONPayloadResponse with status 403
//
// If you do not wish to handle the error, and are ok with a 400 response on error, feel free to use c.JSON(http.StatusForbidden, payload)
func StatusForbiddenJSONPayloadResponse(c *gin.Context, payload any) error {
	return CustomStatusJSONPayloadResponse(c, http.StatusForbidden, payload)
}

// StatusNotFoundResponse is a shorthand for CustomStatusPayloadResponse with status 404
// If you do not wish to handle the error, and are ok with a 400 response on error, feel free to use c.JSON(http.StatusNotFound, payload)
func StatusNotFoundPayloadResponse(c *gin.Context, payload any) {
	CustomStatusPayloadResponse(c, http.StatusNotFound, payload, ContentTypeJSON)
}

// CustomStatusJSONPayloadResponse is a shorthand for CustomStatusPayloadResponse with ContentTypeJSON
//
// If you do not wish to handle the error, and are ok with a 400 response on error, feel free to use c.JSON(status, payload)
func CustomStatusJSONPayloadResponse(c *gin.Context, status int, payload any) error {
	return CustomStatusPayloadResponse(c, status, payload, ContentTypeJSON)
}

// If you do not want to handle the error, and are ok with a 400 response on error, feel free to use gin context's functions
// such as c.JSON, c.XML, etc.
// Valid ContentType: "application/json", "application/xml", "text/html"
//
// "text/html" expects a string as payload
func CustomStatusPayloadResponse(c *gin.Context, status int, payload any, contentType ContentType) error {
	var (
		pL  = []byte{}
		err error
	)
	switch contentType {
	case ContentTypeJSON:
		pL, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	case ContentTypeXML:
		pL, err = xml.Marshal(payload)
		if err != nil {
			return err
		}
	case ContentTypeHTML:
		str, ok := payload.(string)
		if !ok {
			return ErrInvalidPayloadType
		}
		pL = []byte(str)
	default:
		return ErrUnknownContentType
	}
	c.Status(status)
	c.Header("Content-Type", string(contentType))
	_, err = c.Writer.Write(pL)
	return err
}
