package handles

import (
	"net/http"
	"time"

	"github.com/FXAZfung/go-cache"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/model/request"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/internal/service"
	"github.com/FXAZfung/image-board/pkg/utils"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/gin-gonic/gin"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var imageCache = cache.NewMemCache[int]()
var (
	imageDuration = time.Minute
	imageTimes    = 15
)

//// GetImageByID 根据ID获取图片
//// @Summary 根据ID获取图片详情
//// @Description 根据ID获取图片详细信息，包括标签等元数据
//// @Tags 图片
//// @Accept json
//// @Produce json
//// @Param id path int true "图片ID" minimum(1)
//// @Success 200 {object} common.Resp{data=model.Image} "图片详细信息"
//// @Failure 400 {object} common.Resp "ID格式错误"
//// @Failure 404 {object} common.Resp "图片不存在"
//// @Router /images/image/{id} [get]
//func GetImageByID(c *gin.Context) {
//	idStr := c.Param("id")
//	id, err := strconv.ParseUint(idStr, 10, 32)
//	if err != nil {
//		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
//		return
//	}
//
//	image, err := op.GetImageByID(uint(id))
//	if err != nil {
//		common.ErrorResp(c, http.StatusNotFound, err)
//		return
//	}
//
//	common.SuccessResp(c, image)
//}

// GetImageByName 根据文件名获取图片
// @Summary 获取原始图片文件
// @Description 根据文件名直接返回图片二进制内容
// @Tags 图片
// @Produce image/*
// @Param name path string true "文件名" example("example.jpg")
// @Success 200 {file} binary "图片文件"
// @Failure 404 {object} common.Resp "图片不存在"
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

// GetRandomImage 随机获取图片
// @Summary 获取随机图片
// @Description 随机获取一张图片（15分钟内同一IP最多请求15次）
// @Tags 图片
// @Produce image/*
// @Success 200 {file} binary "图片文件"
// @Failure 404 {object} common.Resp "无可用图片"
// @Failure 429 {object} common.Resp "请求过于频繁"
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
// @Summary 分页获取图片列表
// @Description 分页获取所有图片基本信息
// @Tags 图片
// @Accept json
// @Produce json
// @Param page body model.PageReq true "分页参数"
// @Success 200 {object} common.Resp{data=common.PageResp{content=[]model.Image}} "分页结果"
// @Failure 400 {object} common.Resp "参数校验失败"
// @Failure 500 {object} common.Resp "服务器错误"
// @Router /api/image/list [post]
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
// @Summary 上传新图片
// @Description 上传图片文件并添加元数据（需要登录）
// @Tags 图片
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "用户令牌"
// @Param image formData file true "图片文件（支持PNG/JPEG/GIF）"
// @Success 200 {object} common.Resp{data=model.Image} "上传成功"
// @Failure 400 {object} common.Resp "文件无效/参数错误"
// @Failure 401 {object} common.Resp "未授权"
// @Failure 413 {object} common.Resp "文件过大"
// @Failure 500 {object} common.Resp "上传失败"
// @Router /api/image/upload [post]
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
	image, err := service.UploadImage(file, user.(*model.User))
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	// Return success response
	common.SuccessResp(c, image)
}

//// UpdateImage 更新图片信息
//// @Summary 修改图片信息
//// @Description 更新图片描述、可见性等元数据（需要登录）
//// @Tags 认证
//// @Accept json
//// @Produce json
//// @Security ApiKeyAuth
//// @Param Authorization header string true "用户令牌"
//// @Param id path int true "图片ID" minimum(1)
//// @Param image body request.UpdateImageReq true "更新参数"
//// @Success 200 {object} common.Resp{data=model.Image} "更新后的图片信息"
//// @Failure 400 {object} common.Resp "参数错误"
//// @Failure 403 {object} common.Resp "无修改权限"
//// @Failure 404 {object} common.Resp "图片不存在"
//// @Router /api/image/images/{id} [put]
//func UpdateImage(c *gin.Context) {
//	// Parse ID parameter
//	idStr := c.Param("id")
//	id, err := strconv.ParseUint(idStr, 10, 32)
//	if err != nil {
//		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
//		return
//	}
//
//	// Parse request body
//	var req request.UpdateImageReq
//	if err := c.ShouldBind(&req); err != nil {
//		common.ErrorResp(c, http.StatusBadRequest, err)
//		return
//	}
//
//	// Call service to update image
//	image, err := service.UpdateImage(uint(id), req)
//	if err != nil {
//		common.ErrorResp(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	common.SuccessResp(c, image)
//}

// DeleteImage 删除图片
// @Summary 删除图片
// @Description 永久删除图片及其关联数据（需要登录）
// @Tags 图片
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param image body request.ImageDeleteReq true "删除参数"
// @Success 200 {object} common.Resp{data=response.ImageDeleteResponse} "删除结果"
// @Failure 400 {object} common.Resp "参数格式错误"
// @Failure 401 {object} common.Resp "未授权，需要登录"
// @Failure 403 {object} common.Resp "无删除权限"
// @Failure 404 {object} common.Resp "图片不存在"
// @Router /api/image/delete [post]
func DeleteImage(c *gin.Context) {
	var req request.ImageDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	resp, err := service.DeleteImage(req.ID)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, resp)
}

// RemoveTagFromImage 移除图片标签
// @Summary 移除图片标签
// @Description 从图片中移除指定标签（需要登录）
// @Tags 图片
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body request.RemoveTagReq true "移除标签请求"
// @Success 200 {object} common.Resp{data=response.ImageTagResponse} "操作结果"
// @Failure 400 {object} common.Resp "请求参数错误"
// @Failure 401 {object} common.Resp "未授权，需要登录"
// @Failure 404 {object} common.Resp "图片或标签不存在"
// @Failure 500 {object} common.Resp "服务器错误"
// @Router /api/image/tag/remove [post]
func RemoveTagFromImage(c *gin.Context) {
	var req request.RemoveTagReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}
	// Validate ImageID and TagID
	if req.ImageID == 0 || req.TagID == 0 {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	resp, err := service.RemoveTagFromImage(req.ImageID, req.TagID)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, resp)
}

// AddTagToImage 给图片添加标签
// @Summary 添加图片标签
// @Description 给现有图片添加一个标签（需要登录）
// @Tags 图片
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body request.AddTagReq true "添加标签请求"
// @Success 200 {object} common.Resp{data=model.Tag} "添加成功"
// @Failure 400 {object} common.Resp "请求格式无效"
// @Failure 401 {object} common.Resp "未授权，需要登录"
// @Failure 404 {object} common.Resp "图片不存在"
// @Failure 500 {object} common.Resp "服务器错误"
// @Router /api/image/tag/add [post]
func AddTagToImage(c *gin.Context) {

	// Bind request body
	var req request.AddTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	var err error
	req.Tag, err = utils.TrimSpace(req.Tag)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, err)
		return
	}

	// 验证字段
	if req.ID == 0 {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Add tag to image
	tag, err := op.AddTagToImage(req.ID, req.Tag)
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	// Return success response
	common.SuccessResp(c, tag)
}

// GetImagesByTag 按标签搜索图片
// @Summary 根据标签获取图片
// @Description 分页获取包含指定标签的所有图片
// @Tags 图片
// @Accept json
// @Produce json
// @Param tag query string true "标签名称" minLength(1)
// @Param page body model.PageReq true "分页参数"
// @Success 200 {object} common.Resp{data=common.PageResp{content=[]model.Image}} "分页结果"
// @Failure 400 {object} common.Resp "标签参数缺失或格式错误"
// @Failure 404 {object} common.Resp "标签不存在"
// @Failure 500 {object} common.Resp "服务器错误"
// @Router /api/image/tag/list [post]
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

// GetImageCount 获取图片总数
// @Summary 获取图片统计
// @Description 获取系统中的图片总量
// @Tags 图片
// @Produce json
// @Success 200 {object} common.Resp{data=int64} "统计结果"
// @Failure 500 {object} common.Resp "统计失败"
// @Router /api/image/count [get]
func GetImageCount(c *gin.Context) {
	count, err := op.GetImageCount()
	if err != nil {
		common.ErrorResp(c, http.StatusInternalServerError, err)
		return
	}

	common.SuccessResp(c, count)
}

// GetThumbnailByName 获取缩略图
// @Summary 获取图片缩略图
// @Description 获取指定文件的缩略图（自动降级返回原图）
// @Tags 图片
// @Produce image/*
// @Param name path string true "文件名" example("example_thumb.jpg")
// @Success 200 {file} binary "缩略图文件"
// @Failure 404 {object} common.Resp "文件不存在"
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
