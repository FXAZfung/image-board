package model

import "time"

// Tag 标签模型
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null;index"` // 标签名称，唯一且索引
	Count     int       `json:"count" gorm:"default:0"`            // 使用此标签的图片数量
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
