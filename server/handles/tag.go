package handles

import (
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/model/request"
	"github.com/FXAZfung/image-board/internal/model/response"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListTags 列出所有标签（分页）
// @Summary 分页获取标签列表
// @Description 使用分页方式获取所有标签列表
// @Tags 标签
// @Accept json
// @Produce json
// @Param page body model.PageReq true "分页参数"
// @Success 200 {object} common.Resp{data=common.PageResp{content=[]model.Tag}} "分页结果"
// @Failure 400 {object} common.Resp "参数绑定错误"
// @Failure 500 {object} common.Resp "服务器内部错误"
// @Router /api/public/tags [post]
func ListTags(c *gin.Context) {
	var req model.PageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	req.Validate()
	tags, total, err := op.ListTags(req.Page, req.PerPage)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, common.PageResp{
		Content: tags,
		Total:   total,
	})
}

// GetTagByID 通过ID获取标签
// @Summary 根据ID获取标签详情
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} common.Resp{data=model.Tag} "标签详情"
// @Failure 400 {object} common.Resp "ID格式错误"
// @Failure 404 {object} common.Resp "标签不存在"
// @Router /api/public/tags/{id} [get]
func GetTagByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	tag, err := op.GetTagByID(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, err)
		return
	}

	common.SuccessResp(c, tag)
}

// GetTagByName 通过名称获取标签
// @Summary 根据名称查询标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param name query string true "标签名称"
// @Success 200 {object} common.Resp{data=model.Tag} "标签详情"
// @Failure 400 {object} common.Resp "名称参数缺失"
// @Failure 404 {object} common.Resp "标签不存在"
// @Router /api/public/tags/name [get]
func GetTagByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		common.ErrorStrResp(c, http.StatusBadRequest, "Tag name is required")
		return
	}

	tag, err := op.GetTagByName(name)
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, err)
		return
	}

	common.SuccessResp(c, tag)
}

// MostPopularTags 获取最常用标签
// @Summary 获取热门标签列表
// @Description 根据使用次数降序排列获取最常用标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param limit query int false "返回数量限制" minimum(1) default(10)
// @Success 200 {object} common.Resp{data=[]model.Tag} "标签列表"
// @Failure 500 {object} common.Resp "服务器内部错误"
// @Router /api/public/tags/popular [get]
func MostPopularTags(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default to 10 if invalid
	}

	tags, err := op.GetMostPopularTags(limit)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tags)
}

// SearchTags 标签前缀搜索
// @Summary 根据前缀搜索标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param prefix query string true "搜索前缀"
// @Param limit query int false "最大返回数量" minimum(1) default(20)
// @Success 200 {object} common.Resp{data=[]model.Tag} "匹配的标签列表"
// @Failure 400 {object} common.Resp "前缀参数缺失"
// @Failure 500 {object} common.Resp "服务器内部错误"
// @Router /api/public/tags/search [get]
func SearchTags(c *gin.Context) {
	prefix := c.Query("prefix")
	if prefix == "" {
		common.ErrorStrResp(c, http.StatusBadRequest, "Search prefix is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default to 20 if invalid
	}

	tags, err := op.SearchTagsByPrefix(prefix, limit)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tags)
}

// CreateTag 创建新标签
// @Summary 创建新标签
// @Tags 标签
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param tag body request.CreateTagReq true "标签创建参数"
// @Success 200 {object} common.Resp{data=model.Tag} "已创建的标签"
// @Failure 400 {object} common.Resp "参数错误/名称重复"
// @Failure 401 {object} common.Resp "未授权"
// @Failure 500 {object} common.Resp "服务器内部错误"
// @Router /api/auth/tags [post]
func CreateTag(c *gin.Context) {
	var req request.CreateTagReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		common.ErrorStrResp(c, http.StatusBadRequest, "Tag name is required")
		return
	}

	// Check if tag already exists
	_, err := op.GetTagByName(req.Name)
	if err == nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Tag already exists")
		return
	}

	// Create new tag
	tag := &model.Tag{
		Name:  req.Name,
		Count: 0, // Initial count is 0
	}

	if err := op.CreateTag(tag); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tag)
}

// UpdateTag 更新标签
// @Summary 更新标签信息
// @Tags 标签
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "标签ID"
// @Param tag body request.CreateTagReq true "更新后的标签信息"
// @Success 200 {object} common.Resp{data=model.Tag} "更新后的标签"
// @Failure 400 {object} common.Resp "参数错误"
// @Failure 401 {object} common.Resp "未授权"
// @Failure 404 {object} common.Resp "标签不存在"
// @Failure 500 {object} common.Resp "服务器内部错误"
// @Router /api/auth/tags/{id} [put]
func UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req request.CreateTagReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		common.ErrorStrResp(c, http.StatusBadRequest, "Tag name is required")
		return
	}

	// Get existing tag
	tag, err := op.GetTagByID(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, errs.ErrTagNotFound)
		return
	}

	// Update tag name
	tag.Name = req.Name
	if err := op.UpdateTag(tag); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tag)
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Description 删除标签并移除与所有图片的关联
// @Tags 标签
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "标签ID"
// @Success 200 {object} common.Resp{data=response.TagDeleteResponse} "删除结果"
// @Failure 400 {object} common.Resp "ID格式错误"
// @Failure 401 {object} common.Resp "未授权"
// @Failure 404 {object} common.Resp "标签不存在"
// @Failure 500 {object} common.Resp "服务器内部错误"
// @Router /api/auth/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Get tag for response
	tag, err := op.GetTagByID(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, errs.ErrTagNotFound)
		return
	}

	if err := op.DeleteTag(uint(id)); err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, response.TagDeleteResponse{
		ID:      uint(id),
		Name:    tag.Name,
		Success: true,
	})
}

// GetTagsByImage 获取图片标签
// @Summary 获取图片关联标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param image_id path int true "图片ID"
// @Success 200 {object} common.Resp{data=[]model.Tag} "标签列表"
// @Failure 400 {object} common.Resp "ID格式错误"
// @Failure 404 {object} common.Resp "图片不存在"
// @Failure 500 {object} common.Resp "服务器内部错误"
// @Router /api/public/tags/image/{image_id} [get]
func GetTagsByImage(c *gin.Context) {
	imageIDStr := c.Param("image_id")
	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid image ID format")
		return
	}

	// First check if image exists
	_, err = op.GetImageByID(uint(imageID))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, errs.ImageNotFound)
		return
	}

	tags, err := op.GetTagsForImage(uint(imageID))
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, tags)
}
