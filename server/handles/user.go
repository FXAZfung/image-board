package handles

import (
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required" example:"newuser"`     // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
	Role     int    `json:"role" example:"1"`                                  // 角色：2-管理员，0-普通用户 1-游客（只能存在一个）
}

type UpdateUserReq struct {
	Password string `json:"password" example:"newpassword123"` // 新密码
	Role     *int   `json:"role" example:"1"`                  // 角色
	Disable  *bool  `json:"disable" example:"true"`            // 禁用
}

// Register 注册用户
// @Summary 注册新用户
// @Description 创建新用户账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body RegisterReq true "用户信息"
// @Success 200 {object} common.Resp{data=model.User} "注册成功"
// @Failure 403 {object} common.Resp "需要管理权限"
// @Failure 400 {object} common.Resp "参数错误"
// @Failure 409 {object} common.Resp "用户已存在"
// @Router /api/public/register [post]
func Register(c *gin.Context) {
	// 检查权限
	currentUser, _ := c.Get("user")
	if currentUser.(*model.User).Role != model.ADMIN {
		common.ErrorStrResp(c, http.StatusForbidden, "需要管理员权限")
		return
	}
	var req RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	// 创建用户模型
	user := &model.User{
		Username: req.Username,
		Role:     req.Role,
	}

	// 设置密码哈希
	user.SetPassword(req.Password)

	// 保存用户
	if err := op.CreateUser(user); err != nil {
		if err.Error() == errs.ErrUserExist.Error() {
			common.ErrorStrResp(c, http.StatusConflict, "用户名已存在")
			return
		}
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	// 清除密码哈希后返回用户信息
	user.PwdHash = ""
	common.SuccessResp(c, user)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} common.Resp{data=model.User} "用户信息"
// @Failure 401 {object} common.Resp "未授权"
// @Security ApiKeyAuth
// @Router /api/auth/user/info [get]
func GetUserInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		common.ErrorStrResp(c, http.StatusUnauthorized, "未登录")
		return
	}

	// 返回当前用户信息，不包含密码
	userData := user.(*model.User)
	userData.PwdHash = ""
	common.SuccessResp(c, userData)
}

// GetUserByID 根据ID获取用户
// @Summary 根据ID获取用户信息
// @Description 根据用户ID获取用户详细信息（需要管理员权限）
// @Tags 认证
// @Accept json
// @Produce json
// @Param id path int true "用户ID" minimum(1)
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} common.Resp{data=model.User} "用户信息"
// @Failure 400 {object} common.Resp "ID格式错误"
// @Failure 403 {object} common.Resp "权限不足"
// @Failure 404 {object} common.Resp "用户不存在"
// @Security ApiKeyAuth
// @Router /api/auth/users/{id} [get]
func GetUserByID(c *gin.Context) {
	// 检查权限
	currentUser, _ := c.Get("user")
	if currentUser.(*model.User).Role != model.ADMIN {
		common.ErrorStrResp(c, http.StatusForbidden, "需要管理员权限")
		return
	}

	// 解析ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "无效的ID格式")
		return
	}

	// 获取用户
	user, err := op.GetUserById(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, err)
		return
	}

	// 返回用户信息，不包含密码
	user.PwdHash = ""
	common.SuccessResp(c, user)
}

// ListUsers 获取用户列表
// @Summary 分页获取用户列表
// @Description 分页获取所有用户（需要管理员权限）
// @Tags 认证
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1) minimum(1)
// @Param per_page query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} common.Resp{data=common.PageResp{content=[]model.User}} "分页结果"
// @Failure 403 {object} common.Resp "权限不足"
// @Failure 500 {object} common.Resp "服务器错误"
// @Security ApiKeyAuth
// @Router /api/auth/users [get]
func ListUsers(c *gin.Context) {
	// 检查权限
	currentUser, _ := c.Get("user")
	if currentUser.(*model.User).Role != model.ADMIN {
		common.ErrorStrResp(c, http.StatusForbidden, "需要管理员权限")
		return
	}

	var req model.PageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, 400, err)
		return
	}
	req.Validate()

	// 获取用户列表
	users, total, err := op.GetUsers(req.Page, req.PerPage)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	// 清除敏感信息
	for i := range users {
		users[i].PwdHash = ""
	}

	common.SuccessResp(c, common.PageResp{
		Content: users,
		Total:   total,
	})
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 更新指定用户的信息（需要管理员权限或为自己的账号）
// @Tags 认证
// @Accept json
// @Produce json
// @Param id path int true "用户ID" minimum(1)
// @Param user body UpdateUserReq true "用户信息"
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} common.Resp{data=model.User} "更新后的用户信息"
// @Failure 400 {object} common.Resp "参数错误"
// @Failure 403 {object} common.Resp "权限不足"
// @Failure 404 {object} common.Resp "用户不存在"
// @Security ApiKeyAuth
// @Router /api/auth/users/{id} [put]
func UpdateUser(c *gin.Context) {
	// 获取当前用户
	currentUser, _ := c.Get("user")
	currentUserData := currentUser.(*model.User)

	// 解析目标用户ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "无效的ID格式")
		return
	}

	// 检查权限（只能更新自己或管理员更新任何人）
	if currentUserData.ID != uint(id) && currentUserData.Role != model.ADMIN {
		common.ErrorStrResp(c, http.StatusForbidden, "无权限修改其他用户")
		return
	}

	// 获取用户
	user, err := op.GetUserById(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, err)
		return
	}

	// 解析请求
	var req UpdateUserReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	// 更新密码
	if req.Password != "" {
		user.SetPassword(req.Password)
	}

	// 更新角色（只有管理员能更改角色）
	if req.Role != nil && currentUserData.Role == model.ADMIN {
		user.Role = *req.Role
	}

	// 保存更新
	if err := op.UpdateUser(user); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	// 返回更新后的用户，不包含密码
	user.PwdHash = ""
	common.SuccessResp(c, user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定用户（需要管理员权限）
// @Tags 认证
// @Accept json
// @Produce json
// @Param id path int true "用户ID" minimum(1)
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} common.Resp "删除成功"
// @Failure 400 {object} common.Resp "ID格式错误"
// @Failure 403 {object} common.Resp "权限不足"
// @Failure 404 {object} common.Resp "用户不存在"
// @Security ApiKeyAuth
// @Router /api/auth/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	// 检查权限
	currentUser, _ := c.Get("user")
	if currentUser.(*model.User).Role != model.ADMIN {
		common.ErrorStrResp(c, http.StatusForbidden, "需要管理员权限")
		return
	}

	// 解析ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "无效的ID格式")
		return
	}

	// 检查用户是否存在
	_, err = op.GetUserById(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, err)
		return
	}

	// 删除用户
	if err := op.DeleteUser(uint(id)); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, "删除成功")
}

// GetUserCount 获取用户总数
// @Summary 获取用户总数
// @Description 获取系统中的用户总数（需要管理员权限）
// @Tags 认证
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} common.Resp{data=int64} "用户总数"
// @Failure 403 {object} common.Resp "权限不足"
// @Failure 500 {object} common.Resp "服务器错误"
// @Security ApiKeyAuth
// @Router /api/auth/users/count [get]
func GetUserCount(c *gin.Context) {
	// 检查权限
	currentUser, _ := c.Get("user")
	if currentUser.(*model.User).Role != model.ADMIN {
		common.ErrorStrResp(c, http.StatusForbidden, "需要管理员权限")
		return
	}

	// 获取用户总数
	count, err := op.GetUserCount()
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, count)
}
