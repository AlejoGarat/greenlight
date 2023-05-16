package httphelpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReadIDParam(c *gin.Context) (int64, error) {
	params := c.Params

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
