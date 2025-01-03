package model

import (
	"time"
)

type Image struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FileName    string    `json:"file_name" gorm:"unique;not null"`
	ContentType string    `json:"content_type"`
	Path        string    `json:"path" gorm:"not null"`     // 图片路径
	Category    string    `json:"category" gorm:"not null"` // 分类
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	ShortLink   string    `json:"short_link"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}
