package common

import (
	"github.com/FXAZfung/image-board/cmd/flags"
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
)

func hidePrivacy(msg string) string {
	for _, r := range conf.PrivacyReg {
		msg = r.ReplaceAllStringFunc(msg, func(s string) string {
			return strings.Repeat("*", len(s))
		})
	}
	return msg
}

//func ErrorResp(c *gin.Context, code int, err error, l ...bool) {
//	ErrorWithDataResp(c, err, code, nil, l...)
//}

func ErrorWithDataResp(c *gin.Context, code int, data interface{}, err error, l ...bool) {
	if len(l) > 0 && l[0] {
		if flags.Debug || flags.Dev {
			log.Errorf("%+v", err)
		} else {
			log.Errorf("%v", err)
		}
	}
	c.JSON(200, Resp[interface{}]{
		Code:    code,
		Message: hidePrivacy(err.Error()),
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
