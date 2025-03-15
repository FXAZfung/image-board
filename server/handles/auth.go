package handles

import (
	"github.com/FXAZfung/go-cache"
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginReq struct {
	Username string `json:"username" example:"admin" binding:"required"` // 用户名
	Password string `json:"password" example:"admin" binding:"required"` // 密码
}

type LoginResp struct {
	Token  string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT Token
	Expire time.Time   `json:"expire" example:"2023-10-01T12:00:00Z"`                   // Token 过期时间
	User   *model.User `json:"user"`
}

var loginCache = cache.NewMemCache[int]()
var (
	loginDuration = time.Minute * 3
	loginTimes    = 10
)

// Login 登录
// @Summary 登录
// @Description 登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param user body LoginReq true "用户信息"
// @Success 200 {object} LoginResp "登录成功"
// @Failure 400 {object} common.Resp "无效请求"
// @Failure 429 {object} common.Resp "登录尝试次数过多"
// @Failure 500 {object} common.Resp "生成 Token 失败"
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
		common.ErrorResp(c, http.StatusBadRequest, err)
		loginCache.Set(ip, count+1)
		return
	}
	// validate password hash
	if err := user.ValidatePwdStaticHash(req.Password); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		loginCache.Set(ip, count+1)
		return
	}
	// generate token
	token, err := common.GenerateToken(user)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}
	resp := &LoginResp{
		Token:  token,
		Expire: time.Now().Add(time.Duration(conf.Conf.TokenExpiresIn) * time.Hour),
		User:   user,
	}
	common.SuccessResp(c, resp)
	loginCache.Del(ip)
}

// Logout 登出
// @Summary 登出
// @Description 登出
// @Tags 认证
// @Accept json
// @Produce json
// @Param Authorization header string true "Token 格式: {token}" default(<token>)
// @Success 200 {object} common.Resp "登出成功"
// @Failure 500 {object} common.Resp "令牌失效失败"
// @Security ApiKeyAuth
// @Router /api/auth/logout [get]
func Logout(c *gin.Context) {
	err := common.InvalidateToken(c.GetHeader("Authorization"))
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResp(c, "Logout successfully")
	}
}
