package response

import (
	"github.com/FXAZfung/image-board/internal/model"
	"time"
)

// ImageUploadResponse defines the upload image response format
type ImageUploadResponse struct {
	ID            uint   `json:"id" example:"1"`
	OriginalName  string `json:"original_name" example:"abc"`
	FileName      string `json:"file_name" example:"abc"`
	ThumbnailPath string `json:"thumbnail_path" example:"abc123.jpg"`
	WebpPath      string `json:"webp_path" example:"abc123.webp"`
	Path          string `json:"path" example:"abc123.jpg"`
}

type ImageResponse struct {
	ID            uint        `json:"id" example:"1"`
	OriginalName  string      `json:"original_name" example:"abc"`
	FileName      string      `json:"file_name" example:"abc"`
	ThumbnailPath string      `json:"thumbnail_path" example:"abc123.jpg"`
	WebpPath      string      `json:"webp_path" example:"abc123.webp"`
	Path          string      `json:"path" example:"abc123.jpg"`
	Width         int         `json:"width" example:"320"`
	Height        int         `json:"height" example:"240"`
	Size          int64       `json:"size" example:"768"`
	ContentType   string      `json:"content_type" example:"image/jpeg"`
	Description   string      `json:"description" example:"abc123.jpg"`
	IsPublic      bool        `json:"is_public" example:"true"`
	CreatedAt     time.Time   `json:"created_at" example:"2020-01-01 01:01:01"`
	UpdatedAt     time.Time   `json:"updated_at" example:"2020-01-01 01:01:01"`
	Tags          []model.Tag `json:"tags"`
	UserID        uint        `json:"user_id" example:"1"`
	ViewCount     int         `json:"view_count" example:"100"`
	DownloadCount int         `json:"download_count" example:"50"`
}

// ImageDeleteResponse defines the image deletion response format
type ImageDeleteResponse struct {
	Message string `json:"message" example:"Image deleted successfully"`
	ID      uint   `json:"id" example:"1"`
}

// ImageCountResponse defines the response for image count
type ImageCountResponse struct {
	Count int64 `json:"count" example:"42"`
}
