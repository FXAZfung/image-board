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
		common.ErrorResp(c, 400, err)
		return
	}
	req.Validate()
	log.Debugf("%+v", req)
	users, total, err := op.GetUsers(req.Page, req.PerPage)
	if err != nil {
		common.ErrorResp(c, 500, err)
		return
	}
	common.SuccessResp(c, common.PageResp{
		Content: users,
		Total:   total,
	})
}

func CreateUser(c *gin.Context) {
	var req model.User
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, 400, err)
		return
	}
	if req.IsGuest() || req.IsAdmin() {
		common.ErrorStrResp(c, 400, "Can not create guest or admin user")
		return
	}
	err := op.CreateUser(&req)
	if err != nil {
		common.ErrorResp(c, 500, err)
		return
	}
	common.SuccessResp(c, req)
}
