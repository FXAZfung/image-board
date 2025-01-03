package common

import (
	"github.com/gin-gonic/gin"
)

func ErrorWithDataResp(c *gin.Context, code int, data interface{}, err error) {
	c.JSON(200, Resp[interface{}]{
		Code:    code,
		Message: err.Error(),
		Data:    data,
	})
	c.Abort()
}

func ErrorStrResp(c *gin.Context, code int, str string) {
	c.JSON(200, Resp[interface{}]{
		Code:    code,
		Message: str,
		Data:    nil,
	})
	c.Abort()
}

func SuccessResp(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, Resp[interface{}]{
			Code:    200,
			Message: "success",
			Data:    nil,
		})
		return
	}
	c.JSON(200, Resp[interface{}]{
		Code:    200,
		Message: "success",
		Data:    data[0],
	})
}
