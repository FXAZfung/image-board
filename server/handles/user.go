package handles

import (
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ListUser 列出用户列表
// @Summary 列出用户列表
// @Description 列出用户列表
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param page body model.PageReq true "分类信息"
// @Success 200 {object} common.PageResp
// @Router /api/private/users [post]
func ListUser(c *gin.Context) {
	var req model.PageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorStrResp(c, 400, err.Error())
		return
	}
	req.Validate()
	log.Debugf("%+v", req)
	users, total, err := op.GetUsers(req.Page, req.PerPage)
	if err != nil {
		common.ErrorStrResp(c, 500, err.Error())
		return
	}
	common.SuccessResp(c, common.PageResp{
		Content: users,
		Total:   total,
	})
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建用户
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param user body model.User true "用户信息"
// @Success 200 {object} common.EmptyResp
// @Router /api/private/user [post]
func CreateUser(c *gin.Context) {

}
