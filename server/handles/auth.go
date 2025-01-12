package handles

import (
	"github.com/FXAZfung/go-cache"
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

var loginCache = cache.NewMemCache[int]()
var (
	loginDuration = time.Minute * 3
	loginTimes    = 10
)

// Login 登录
// @Summary 登录
// @Description 登录
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginReq true "用户信息"
// @Success 200 {object} LoginResp "登录成功"
// @Router /api/public/login [post]
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Bad request")
		return
	}
	loginHash(c, &req)
}

func loginHash(c *gin.Context, req *LoginReq) {
	// check count of login
	ip := c.ClientIP()
	count, ok := loginCache.Get(ip)
	if ok && count >= loginTimes {
		common.ErrorStrResp(c, http.StatusTooManyRequests, "Too many unsuccessful sign-in attempts have been made using an incorrect username or password, Try again later.")
		loginCache.Expire(ip, loginDuration)
		return
	}
	// check username
	user, err := op.GetUserByName(req.Username)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, err.Error())
		loginCache.Set(ip, count+1)
		return
	}
	// validate password hash
	if err := user.ValidatePwdStaticHash(req.Password); err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, err.Error())
		loginCache.Set(ip, count+1)
		return
	}
	// generate token
	token, err := common.GenerateToken(user)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	resp := &LoginResp{
		Token:  token,
		Expire: time.Now().Add(time.Duration(conf.Conf.TokenExpiresIn) * time.Hour),
	}
	common.SuccessResp(c, resp)
	loginCache.Del(ip)
}

// Logout 登出
// @Summary 登出
// @Description 登出
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Token"
// @Success 200 {string} "登出成功"
// @Router /api/auth/logout [get]
func Logout(c *gin.Context) {
	err := common.InvalidateToken(c.GetHeader("Authorization"))
	if err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		common.SuccessResp(c)
	}
}
