package middleware

import (
	"crypto/subtle"
	"github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/internal/setting"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if subtle.ConstantTimeCompare([]byte(token), []byte(setting.GetStr(config.Token))) == 1 {
		admin, err := op.GetAdmin()
		if err != nil {
			common.ErrorStrResp(c, http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		c.Set("user", admin)
		log.Debugf("use admin token: %+v", admin)
		c.Next()
		return
	}
	if token == "" {
		guest, err := op.GetGuest()
		if err != nil {
			common.ErrorStrResp(c, http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		if guest.Disabled {
			common.ErrorStrResp(c, http.StatusUnauthorized, "Guest user is disabled, login please")
			c.Abort()
			return
		}
		c.Set("user", guest)
		log.Debugf("use empty token: %+v", guest)
		c.Next()
		return
	}
	userClaims, err := common.ParseToken(token)
	if err != nil {
		common.ErrorStrResp(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}
	user, err := op.GetUserByName(userClaims.Username)
	if err != nil {
		common.ErrorStrResp(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}
	// validate password timestamp
	//if userClaims.PwdTS != user.PwdTS {
	//	common.ErrorStrResp(c, "Password has been changed, login please", 401)
	//	c.Abort()
	//	return
	//}
	if user.Disabled {
		common.ErrorStrResp(c, http.StatusUnauthorized, "Current user is disabled, replace please")
		c.Abort()
		return
	}
	c.Set("user", user)
	log.Debugf("use login token: %+v", user)
	c.Next()
}
