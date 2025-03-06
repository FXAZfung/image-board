package model

import (
	"time"
)

// Image 图片模型
type Image struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	FileName      string    `json:"file_name" gorm:"unique;not null"`
	OriginalName  string    `json:"original_name"`
	Hash          string    `json:"hash" gorm:"unique;not null"`
	Path          string    `json:"path" gorm:"-:all"`           // 图片路径
	ThumbnailPath string    `json:"thumbnail_path" gorm:"-:all"` // 缩略图路径
	ContentType   string    `json:"content_type"`
	Size          int64     `json:"size"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	Description   string    `json:"description"`
	IsPublic      bool      `json:"is_public" gorm:"default:true"`
	ViewCount     int       `json:"view_count" gorm:"default:0"`     // 浏览次数
	DownloadCount int       `json:"download_count" gorm:"default:0"` // 下载次数
	IsDeleted     bool      `json:"is_deleted" gorm:"default:false"`
	UserID        uint      `json:"user_id"`
	Tags          []Tag     `json:"tags" gorm:"many2many:im_image_tags;"` // 标签，多对多关系
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ImageTag 图片与标签的关联表
type ImageTag struct {
	ImageID   uint      `gorm:"primaryKey"`
	TagID     uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
