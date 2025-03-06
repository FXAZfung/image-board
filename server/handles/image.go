package handles

import (
	"github.com/FXAZfung/go-cache"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/model/request"
	"github.com/FXAZfung/image-board/internal/model/response"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/internal/service"
	"github.com/FXAZfung/image-board/pkg/utils"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var imageCache = cache.NewMemCache[int]()
var (
	imageDuration = time.Minute
	imageTimes    = 15
)

// GetImageByID 根据ID获取图片
// @Summary 根据ID获取图片
// @Description 根据ID获取图片详细信息，包括标签等
// @Tags image
// @Accept json
// @Produce json
// @Param id path int true "图片ID"
// @Success 200 {object} model.Image "图片信息"
// @Router /api/public/images/{id} [get]
func GetImageByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	image, err := op.GetImageByID(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusNotFound, err)
		return
	}

	common.SuccessResp(c, image)
}

// GetImageByName 根据文件名获取图片
// @Summary 根据文件名获取图片
// @Description 直接返回图片文件
// @Tags image
// @Accept json
// @Produce image/*
// @Param name path string true "文件名"
// @Success 200 {file} binary "图片内容"
// @Router /images/image/{name} [get]
func GetImageByName(c *gin.Context) {
	name := c.Param("name")

	imageData, err := op.GetImageByFileName(name)
	if err != nil || imageData == nil {
		common.ErrorStrResp(c, http.StatusNotFound, "Image not found")
		return
	}
	c.File(imageData.Path)
}

// GetRandomImage 随机获取一个图片
// @Summary 随机获取一个图片
// @Description 随机获取一个图片，支持按分类过滤
// @Tags image
// @Accept json
// @Produce image/*
// @Success 200 {file} binary "图片内容"
// @Router /images/image/random [get]
func GetRandomImage(c *gin.Context) {
	// 检查请求频率限制
	ip := c.ClientIP()
	count, ok := imageCache.Get(ip)
	if ok && count >= imageTimes {
		common.ErrorStrResp(c, http.StatusTooManyRequests, "Too many requests for random images")
		imageCache.Expire(ip, imageDuration)
		return
	}

	// 设置不缓存头
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	imageData, err := op.GetRandomImage()
	if err != nil || imageData == nil {
		common.ErrorStrResp(c, http.StatusNotFound, "Image not found")
		return
	}

	// 更新请求计数器
	if ok {
		imageCache.Set(ip, count+1)
	} else {
		imageCache.Set(ip, 1, cache.WithEx[int](imageDuration))
	}

	c.File(imageData.Path)
}

// ListImages 分页列出图片
// @Summary 分页列出图片
// @Description 分页获取所有图片
// @Tags image
// @Accept json
// @Produce json
// @Param page body model.PageReq true "分页参数"
// @Success 200 {object} common.PageResp{content=model.Image} "图片列表和总数"
// @Router /api/public/images [post]
func ListImages(c *gin.Context) {
	var req model.PageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	req.Validate()
	images, total, err := op.GetImagesByPage(req.Page, req.PerPage)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, common.PageResp{
		Content: images,
		Total:   total,
	})
}

// UploadImage 上传图片
// @Summary 上传图片
// @Description 上传新图片并可选添加描述、主标签等信息
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param image formData file true "图片文件"
// @Param description formData string false "描述"
// @Param is_public formData boolean false "是否公开"
// @Param main_tag formData string false "主标签"
// @Success 200 {object} response.ImageUploadResponse "上传图片成功响应"
// @Router /api/auth/upload [post]
func UploadImage(c *gin.Context) {
	// Parse request
	var req request.UploadImageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	// Get image file
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Missing image file")
		return
	}
	req.Image = file

	// Get current user
	user, exist := c.Get("user")
	if !exist {
		common.ErrorStrResp(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to upload image
	image, err := service.UploadImage(file, user.(*model.User), req)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	// Return success response
	common.SuccessResp(c, response.ImageUploadResponse{
		ID:   image.ID,
		Path: image.FileName,
	})
}

// UpdateImage 更新图片信息
// @Summary 更新图片信息
// @Description 更新图片的描述、可见性等信息
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "图片ID"
// @Param image body request.UpdateImageReq true "更新信息"
// @Success 200 {object} model.Image "图片更新成功相应"
// @Router /api/auth/images/{id} [put]
func UpdateImage(c *gin.Context) {
	// Parse ID parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Parse request body
	var req request.UpdateImageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	// Call service to update image
	image, err := service.UpdateImage(uint(id), req)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, image)
}

// DeleteImage 删除图片
// @Summary 删除图片
// @Description 删除指定ID的图片及其关联数据
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "图片ID"
// @Success 200 {object} response.ImageDeleteResponse "图片删除成功响应"
// @Router /api/auth/images/{id} [delete]
func DeleteImage(c *gin.Context) {
	// Parse ID parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Call service to delete image
	resp, err := service.DeleteImage(uint(id))
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, resp)
}

// RemoveTagFromImage 从图片中移除标签
// @Summary 从图片中移除标签
// @Description 从图片中移除指定标签
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "图片ID"
// @Param tag_id path int true "标签ID"
// @Success 200 {object} response.ImageTagResponse "移除标签成功响应"
// @Router /api/auth/images/{id}/tags/{tag_id} [delete]
func RemoveTagFromImage(c *gin.Context) {
	// Parse parameters
	imageIDStr := c.Param("id")
	tagIDStr := c.Param("tag_id")

	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid image ID format")
		return
	}

	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid tag ID format")
		return
	}

	// Call service to remove tag
	resp, err := service.RemoveTagFromImage(uint(imageID), uint(tagID))
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, resp)
}

// AddTagsToImage 给图片添加标签
// @Summary 给图片添加标签
// @Description 给图片添加一个或多个标签
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "图片ID"
// @Param tags body request.AddTagsReq true "标签列表"
// @Success 200 {object} response.ImageTagResponse "添加标签成功响应"
// @Router /api/auth/images/{id}/tags [post]
func AddTagsToImage(c *gin.Context) {
	// Parse ID parameter
	imageIDStr := c.Param("id")
	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid image ID format")
		return
	}

	// Parse request body
	var req request.AddTagsReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	if len(req.Tags) == 0 {
		common.ErrorStrResp(c, http.StatusBadRequest, "No tags provided")
		return
	}

	// Call service to add tags
	resp, err := service.AddTagsToImage(uint(imageID), req.Tags)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, resp)
}

// GetImagesByTag handles requests to retrieve images by a specific tag name
// @Summary 根据标签获取图片
// @Description 获取包含特定标签的所有图片
// @Tags image
// @Accept json
// @Produce json
// @Param tag query string true "标签名称"
// @Param page body model.PageReq true "分页参数"
// @Success 200 {object} common.PageResp{content=model.Image} "图片列表和总数"
// @Router /api/public/images/tag [post]
func GetImagesByTag(c *gin.Context) {
	tagName := c.Query("tag")
	if tagName == "" {
		common.ErrorStrResp(c, 400, "Tag name is required")
		return
	}

	var req model.PageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, 400, err)
		return
	}
	req.Validate()

	images, count, err := op.GetImagesByTag(tagName, req.Page, req.PerPage)
	if err != nil {
		common.ErrorResp(c, 500, err)
		return
	}

	common.SuccessResp(c, common.PageResp{
		Content: images,
		Total:   count,
	})
}

// GetImageCount 获取图片数量
// @Summary 获取图片数量
// @Description 获取系统中的图片总数
// @Tags image
// @Accept json
// @Produce json
// @Success 200 {object} response.ImageCountResponse "图片数量"
// @Router /api/public/images/count [get]
func GetImageCount(c *gin.Context) {
	count, err := op.GetImageCount()
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, response.ImageCountResponse{
		Count: count,
	})
}

// GetThumbnailByName 获取图片缩略图
// @Summary 获取图片缩略图
// @Description 根据文件名获取图片的缩略图
// @Tags image
// @Accept json
// @Produce image/*
// @Param name path string true "文件名"
// @Success 200 {file} binary "缩略图内容"
// @Router /images/thumbnail/{name} [get]
func GetThumbnailByName(c *gin.Context) {
	name := c.Param("name")

	imageData, err := op.GetImageByFileName(name)
	if err != nil || imageData == nil {
		common.ErrorStrResp(c, http.StatusNotFound, "Image not found")
		return
	}

	// 构建缩略图路径
	thumbnailPath := service.GetThumbnailPath(imageData.Path)

	// 如果缩略图不存在，则重定向到原图
	if !utils.IsExist(thumbnailPath) {
		c.File(imageData.Path)
		return
	}

	c.File(thumbnailPath)
}
