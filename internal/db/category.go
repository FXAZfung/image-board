package db

import (
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/pkg/errors"
)

// GetCategoryByName 根据分类名获取分类
func GetCategoryByName(name string) (*model.Category, error) {
	var category model.Category
	err := db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, errors.WithStack(errs.ErrCategoryGet)
	}
	return &category, nil
}

// GetCategories 获取所有分类
func GetCategories() ([]*model.Category, error) {
	var categories []*model.Category
	err := db.Find(&categories).Error
	if err != nil {
		return nil, errors.WithStack(errs.ErrCategoryGet)
	}
	return categories, nil
}

// CreateCategory 创建分类
func CreateCategory(category *model.Category) error {
	err := db.Create(category).Error
	if err != nil {
		return errors.WithStack(errs.ErrCategoryCreate)
	}
	return nil
}

// UpdateCategory 更新分类
func UpdateCategory(category *model.Category) error {
	err := db.Save(category).Error
	if err != nil {
		return errors.WithStack(errs.ErrCategoryUpdate)
	}
	return nil
}

// DeleteCategory 删除分类
func DeleteCategory(category *model.Category) error {
	err := db.Delete(category).Error
	if err != nil {
		return errors.WithStack(errs.ErrCategoryDelete)
	}
	return nil
}

// GetCategoryCount 获取分类数量
func GetCategoryCount() (int64, error) {
	var count int64
	err := db.Model(&model.Category{}).Count(&count).Error
	if err != nil {
		return 0, errors.WithStack(errs.ErrCategoryCount)
	}
	return count, nil
}
