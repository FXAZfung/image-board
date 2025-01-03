package model

type Category struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"unique;not null"`
	IsRandom bool   `json:"is_random"`
	IsPublic bool   `json:"is_public"`
}
