package handles

import (
	"github.com/FXAZfung/image-board/cmd/flags"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InfoResp struct {
	ImageCount    int64 `json:"image_count"`
	CategoryCount int64 `json:"category_count"`
	StorageUsage  int64 `json:"storage_usage"`
	UserCount     int64 `json:"user_count"`
}

// GetInfo 获取总图片数，分类数，存储空间使用情况，用户数
// @Summary 获取信息
// @Description 获取信息
// @Tags info
// @Accept json
// @Produce json
// @Success 200 {object} InfoResp "信息"
// @Router /api/public/info [get]
func GetInfo(c *gin.Context) {
	// 获取总图片数
	imageCount, err := op.GetImageCount()
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}
	// 获取分类数
	categoryCount, err := op.GetCategoryCount()
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}
	// 获取存储空间使用情况
	storage, err := op.GetStorageUsage(flags.DataDir)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}
	// 获取用户数
	userCount, err := op.GetUserCount()
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}
	common.SuccessResp(c, InfoResp{
		ImageCount:    imageCount,
		CategoryCount: categoryCount,
		StorageUsage:  storage,
		UserCount:     userCount,
	})
}

//// GetOperation 查询最近的操作
//// @Summary 查询最近的操作
//// @Description 查询最近的操作
//// @Tags info
//// @Accept json
//// @Produce json
//// @Param Authorization header string true "Token"
//// @Param limit query string false "限制数量"
//// @Router /api/private/operation [get]
//func GetOperation(c *gin.Context) {
//	limit := c.DefaultQuery("limit", "10")
//	ops, err := op.GetOperation(limit)
//	if err != nil {
//		common.ErrorStrResp(c, http.StatusInternalServerError, err.Error())
//		return
//	}
//	common.SuccessResp(c, ops)
//}
