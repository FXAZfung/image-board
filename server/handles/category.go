package handles

import (
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryReq struct {
	Name     string `json:"name"`
	IsRandom bool   `json:"is_random"`
	IsPublic bool   `json:"is_public"`
}

// GetCategories 获取图片分类
// @Summary 获取图片分类
// @Description 获取图片分类
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} []model.Category "分类列表"
// @Router /api/public/categories [get]
func GetCategories(c *gin.Context) {
	categories, err := op.GetCategories()
	if err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}
	common.SuccessResp(c, categories)
}

// GetCategoryByName 根据分类名获取图片
// @Summary 根据分类名获取图片
// @Description 根据分类名获取图片
// @Tags category
// @Accept json
// @Produce json
// @Param name path string true "分类名"
// @Success 200 {object} model.Category "分类"
// @Router /api/public/categories/{name} [get]
func GetCategoryByName(c *gin.Context) {
	name := c.Param("name")

	category, err := op.GetCategoryByName(name)
	if err != nil {
		common.ErrorStrResp(c, http.StatusNotFound, "Category not found")
		return
	}
	common.SuccessResp(c, category)
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Description 创建分类
// @Tags category
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @Param category body CategoryReq true "分类信息"
// @Success 200 {object} model.Category "分类"
// @Router /api/auth/categories [post]
func CreateCategory(c *gin.Context) {
	var req CategoryReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Bad request")
		return
	}
	item := &model.Category{
		Name:     req.Name,
		IsRandom: req.IsRandom,
		IsPublic: req.IsPublic,
	}
	err := op.SaveCategory(item)
	if err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to create category")
		return
	}
	common.SuccessResp(c, item)
}
