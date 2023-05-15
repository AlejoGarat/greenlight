package httphelpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const maxBytes int64 = 1_048_576

// JSONDecode will try to decode json into pointer v. In case of unknown fields, they will be ignored
//
// If you do not wish to handle the error and are fine with 400 response on error
// feel free to use c.BindJSON(v)
func JSONDecode(c *gin.Context, v any) error {
	return jsonDecode(c, v, true)
}

// JSONDecode will try to decode json into pointer v. In case of unknown fields, an error will be returned
//
// If you do not wish to handle the error and are fine with 400 response on error
// feel free to use c.BindJSON(v)
func JSONDecodeNoUnknownFieldsAllowed(c *gin.Context, v any) error {
	return jsonDecode(c, v, false)
}

func jsonDecode(c *gin.Context, v any, allowUnknownFields bool) error {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

	decoder := json.NewDecoder(c.Request.Body)
	if !allowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	err := decoder.Decode(v)
	if err != nil {
		return err
	}

	if e := decoder.Decode(&struct{}{}); e != io.EOF {
		err = errors.New("body must only contain a single JSON value")
		return err
	}

	return nil
}
