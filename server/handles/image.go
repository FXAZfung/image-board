package handles

import (
	"bytes"
	"github.com/FXAZfung/go-cache"
	"github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/pkg/random"
	"github.com/FXAZfung/image-board/server/common"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var imageCache = cache.NewMemCache[int]()
var (
	imageDuration = time.Minute
	imageTimes    = 15
)

type ImageDeleteReq struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// GetImageByName 根据文件名获取图片
//
//	@Summary		根据文件名获取图片
//	@title		    根据文件名获取图片
//	@version		1.0
//	@Description	根据文件名获取图片
//	@termsOfService	http://www.swagger.io/terms/
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"文件名"
//	@Success		200		{string}	string	"图片内容"
//	@Router			/api/public/images/{name} [get]
func GetImageByName(c *gin.Context) {
	name := c.Param("name")

	imageData, err := op.GetImageByFileName(name)
	if err != nil || imageData == nil {
		common.ErrorStrResp(c, http.StatusNotFound, "Image not found")
		return
	}
	c.File(imageData.Path)
}

// GetRandomImage 随机获取一个图片 支持分类
// @Summary 随机获取一个图片 支持分类
// @Description 随机获取一个图片 支持分类
// @Tags image
// @Accept json
// @Produce json
// @Param category query string false "分类"
// @Success 200 {object} string "图片内容"
// @Router /api/public/random [get]
func GetRandomImage(c *gin.Context) {

	// check count of login
	ip := c.ClientIP()
	count, ok := imageCache.Get(ip)
	if ok && count >= imageTimes {
		common.ErrorStrResp(c, http.StatusTooManyRequests, "Too many requests for image in a short time")
		imageCache.Expire(ip, imageDuration)
		return
	}
	category := c.Query("category")

	c.Header("Cache-Control", "no-cache")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	imageData, err := op.GetRandomImage(category)
	if err != nil || imageData == nil {
		common.ErrorStrResp(c, http.StatusNotFound, "Image not found")
		return
	}
	imageCache.Set(ip, count+1)
	c.File(imageData.Path)
}

// ListImages 分页列出图片
// @Summary 分页列出图片
// @Description 分页列出图片
// @Tags image
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} string "图片列表"
// @Router /api/public/images [get]
func ListImages(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid page size")
		return
	}

	images, err := op.GetImagesByPage(page, pageSize)
	if err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to retrieve images")
		return
	}

	common.SuccessResp(c, images)
}

// GetImageByShortLink 根据短链获取图片
// @Summary 根据短链获取图片
// @Description 根据短链获取图片
// @Tags image
// @Accept json
// @Produce json
// @Param short_link path string true "短链"
// @Success 200 {object} string "图片内容"
// @Router /api/public/short/{short_link} [get]
func GetImageByShortLink(c *gin.Context) {
	shortLink := c.Param("short_link")

	imageData, err := op.GetImageByShortLink(shortLink)
	if err != nil || imageData == nil {
		common.ErrorStrResp(c, http.StatusNotFound, "Image not found")
		return
	}

	c.File(imageData.Path)
}

//TODO
//// ReloadImages 重新从磁盘加载图片同时更新数据库
//// @Summary 重新从磁盘加载图片同时更新数据库
//// @Description 重新从磁盘加载图片同时更新数据库
//// @Tags image
//// @Accept json
//// @Produce json
//// @Param Authorization header string true "Token"
//// @Success 200 {string} "成功"
//// @Router /api/private/reload [get]
//func ReloadImages(c *gin.Context) {
//	err := op.ReloadImages()
//	if err != nil {
//		common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to reload images")
//		return
//	}
//	common.SuccessResp(c)
//}

// UploadImage 上传图片
// @Summary 上传图片
// @Description 上传图片
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Token"
// @Param image formData file true "图片"
// @Param short_link formData string false "自定义短链"
// @Param category formData string false "分类"
// @Success 200 {object} string "短链"
// @Router /api/auth/upload [post]
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Failed to get file")
		return
	}

	// 读取文件内容
	f, err := file.Open()
	if err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to read file")
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to read file content")
		return
	}

	// 获取图片宽高
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Invalid image format "+err.Error())
		return
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// 对上传的图片的文件名进行处理
	file.Filename = random.RandomizeFileName(file.Filename)

	// 检查是否有自定义短链
	customShortLink := c.PostForm("short_link")
	if customShortLink == "" {
		customShortLink = random.String(6)
	} else {
		// 检查短链是否唯一
		existing, _ := op.GetImageByShortLink(customShortLink)
		if existing != nil {
			common.ErrorStrResp(c, http.StatusConflict, "Custom short link already exists")
			return
		}
	}

	// 统一使用小写分类
	category := strings.ToLower(c.PostForm("category"))

	targetFilePath := path.Join(config.Conf.DataImage.Dir, category, file.Filename)

	imageData := &model.Image{
		FileName:    file.Filename,
		Width:       width,
		Height:      height,
		ShortLink:   customShortLink,
		Path:        targetFilePath,
		ContentType: file.Header.Get("Content-Type"),
		Category:    category,
	}
	// 保存图片到数据库
	err = op.CreateImage(imageData, data, targetFilePath)
	if err != nil {
		common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to save image")
		return
	}

	common.SuccessResp(c, customShortLink)
}

// UploadImages 批量上传图片
// @Summary 批量上传图片
// @Description 批量上传图片
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Token"
// @Param images formData file true "图片"
// @Param category formData string false "分类"
// @Router /api/auth/upload [post]
func UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		common.ErrorStrResp(c, http.StatusBadRequest, "Failed to get file")
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		common.ErrorStrResp(c, http.StatusBadRequest, "Failed to get file")
		return
	}

	// 限制最大上传文件数为10
	if len(files) > 10 {
		common.ErrorStrResp(c, http.StatusBadRequest, "Maximum 10 files allowed")
		return
	}

	category := c.PostForm("category")

	for _, file := range files {
		// 将文件挨个保存
		f, err := file.Open()
		if err != nil {
			common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to read file")
			return // 有一个文件读取失败就返回
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to read file content")
			return
		}
		// 获取图片宽高
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			common.ErrorStrResp(c, http.StatusBadRequest, "Invalid image format "+err.Error())
			return
		}
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()

		// 对上传的图片的文件名进行处理
		file.Filename = random.RandomizeFileName(file.Filename)

		targetFilePath := path.Join(config.Conf.DataImage.Dir, category, file.Filename)

		imageData := &model.Image{
			FileName:    file.Filename,
			Width:       width,
			Height:      height,
			ShortLink:   random.String(6),
			Path:        targetFilePath,
			ContentType: file.Header.Get("Content-Type"),
			Category:    category,
		}
		// 保存图片到数据库
		err = op.CreateImage(imageData, data, targetFilePath)
		if err != nil {
			common.ErrorStrResp(c, http.StatusInternalServerError, "Failed to save image")
			return
		}
	}
	common.SuccessResp(c, "Uploads successful")
}
