package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetInt64InQuery(c *gin.Context, name string) (int64, error) {
	str, ok := c.GetQuery(name)
	if !ok {
		return 0, errors.New("miss parameter " + name + " InQuery")
	}
	return strconv.ParseInt(str, 10, 64)
}
