package handles

import (
	"net/http"
	"time"

	"github.com/FXAZfung/go-cache"
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
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
// @Summary 用户登录
// @Description 通过用户名和密码登录系统并获取令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param user body LoginReq true "用户登录信息"
// @Success 200 {object} common.Resp{data=LoginResp} "登录成功，返回令牌信息"
// @Failure 400 {object} common.Resp "无效请求或密码错误"
// @Failure 429 {object} common.Resp "登录尝试次数过多，请稍后再试"
// @Failure 500 {object} common.Resp "生成令牌失败"
// @Router /api/auth/login [post]
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
// @Summary 用户登出
// @Description 使当前用户令牌失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} common.Resp "登出成功"
// @Failure 401 {object} common.Resp "未授权，缺少有效令牌"
// @Failure 500 {object} common.Resp "令牌失效操作失败"
// @Router /api/auth/logout [post]
func Logout(c *gin.Context) {
	err := common.InvalidateToken(c.GetHeader("Authorization"))
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	} else {
		common.SuccessResp(c, "Logout successfully")
	}
}
