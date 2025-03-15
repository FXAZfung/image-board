package request

import (
	"mime/multipart"
)

// UploadImageReq 上传图片请求
type UploadImageReq struct {
	Image *multipart.FileHeader `json:"-" form:"image" binding:"required"`
}

// UpdateImageReq 更新图片请求
type UpdateImageReq struct {
	Description string `json:"description" form:"description"`
	IsPublic    *bool  `json:"is_public" form:"is_public"`
}

// ImageSearchReq 图片搜索请求
type ImageSearchReq struct {
	Tags      []string `json:"tags"`
	StartDate string   `json:"start_date"`
	EndDate   string   `json:"end_date"`
}
