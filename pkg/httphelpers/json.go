package httphelpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const maxBytes int64 = 1_048_576

type Envelope map[string]any

func JSONDecode(c *gin.Context, v any) error {
	return jsonDecode(c, v, true)
}

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

func WriteJSON(c *gin.Context, status int, data Envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		c.Writer.Header()[key] = value
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	c.Writer.Write(js)

	return nil
}
