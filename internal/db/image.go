package db

import (
	"fmt"
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateImage 创建图片并关联主标签
func CreateImage(image *model.Image) error {
	return db.Create(image).Error
}

// CreateImageWithTags creates an image with associated tags
func CreateImageWithTags(image *model.Image, tagNames []string) error {
	// Start transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// Create image
		if err := tx.Create(image).Error; err != nil {
			return err
		}

		// Add tags to image
		if len(tagNames) > 0 {
			return AddTagsToImage(image.ID, tagNames)
		}

		return nil
	})
}

// GetImageByID retrieves an image by ID
func GetImageByID(id uint) (*model.Image, error) {
	var image model.Image
	err := db.Preload("Tags").First(&image, id).Error
	if err != nil {
		return nil, errors.WithStack(errs.ImageNotFound)
	}
	return &image, nil
}

// GetImageByFilename 根据文件名获取图片
func GetImageByFilename(filename string) (*model.Image, error) {
	var image model.Image
	if err := db.Where("file_name = ?", filename).First(&image).Error; err != nil {
		return nil, errors.WithStack(errs.ImageNotFound)
	}
	return &image, nil
}

// GetImageByHash 根据hash获取图片
func GetImageByHash(hash string) (*model.Image, error) {
	var image model.Image
	if err := db.Where("hash = ?", hash).First(&image).Error; err != nil {
		return nil, errors.WithStack(errs.ImageNotFound)
	}
	return &image, nil
}

// GetImagesByPage 分页获取所有图片
func GetImagesByPage(page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	var count int64
	if err := db.Model(&model.Image{}).Count(&count).Error; err != nil {
		return nil, 0, errors.WithStack(errs.ErrImageCount)
	}
	if err := db.Preload("MainTag").Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error; err != nil {
		return nil, 0, errors.WithStack(errs.ErrImageList)
	}
	return images, count, nil
}

// GetImagesByTag retrieves images with a specific tag
func GetImagesByTag(tagName string, page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	var count int64

	// First verify the tag exists
	var tag model.Tag
	if err := db.Where("name = ?", tagName).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Tag doesn't exist, return empty result
			return []*model.Image{}, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to find tag: %w", err)
	}

	// Count images
	countErr := db.Model(&model.Image{}).
		Joins("INNER JOIN im_image_tags ON im_image_tags.image_id = im_images.id").
		Where("im_image_tags.tag_id = ?", tag.ID).
		Count(&count).Error

	if countErr != nil {
		log.WithFields(log.Fields{
			"tag_id":   tag.ID,
			"tag_name": tagName,
			"error":    countErr,
		}).Error("Database error counting images by tag")
		return nil, 0, errors.WithStack(errs.ErrImageCount)
	}

	// If no images found, return empty result
	if count == 0 {
		return []*model.Image{}, 0, nil
	}

	// Get images
	queryErr := db.Preload("Tags").
		Table("im_images").
		Joins("INNER JOIN im_image_tags ON im_image_tags.image_id = im_images.id").
		Where("im_image_tags.tag_id = ?", tag.ID).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&images).Error

	if queryErr != nil {
		log.WithFields(log.Fields{
			"tag_id":   tag.ID,
			"tag_name": tagName,
			"error":    queryErr,
		}).Error("Database error retrieving images by tag")
		return nil, 0, errors.WithStack(errs.ErrImageList)
	}

	return images, count, nil
}

// UpdateImage 更新图片信息
func UpdateImage(image *model.Image) error {
	return db.Model(image).Updates(map[string]interface{}{
		"description": image.Description,
		"is_public":   image.IsPublic,
		"width":       image.Width,
		"height":      image.Height,
	}).Error
}

// DeleteImage 删除图片及其标签关联
func DeleteImage(imageID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var image model.Image
		if err := tx.Preload("Tags").First(&image, imageID).Error; err != nil {
			return err
		}

		// 删除与标签的关联并减少标签计数
		for _, tag := range image.Tags {
			if err := tx.Where("image_id = ? AND tag_id = ?", imageID, tag.ID).Delete(&model.ImageTag{}).Error; err != nil {
				return err
			}
			if err := tx.Model(&tag).Update("count", tag.Count-1).Error; err != nil {
				return err
			}
		}

		// 删除图片记录
		return tx.Delete(&image).Error
	})
}

// RemoveTagFromImage 从图片中移除标签
func RemoveTagFromImage(imageID uint, tagID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 删除关联
		if err := tx.Where("image_id = ? AND tag_id = ?", imageID, tagID).Delete(&model.ImageTag{}).Error; err != nil {
			return err
		}

		// 减少标签计数
		var tag model.Tag
		if err := tx.First(&tag, tagID).Error; err != nil {
			return err
		}
		return tx.Model(&tag).Update("count", tag.Count-1).Error
	})
}

// GetImageCount 获取图片总数
func GetImageCount() (int64, error) {
	var count int64
	if err := db.Model(&model.Image{}).Count(&count).Error; err != nil {
		return 0, errors.WithStack(errs.ErrImageCount)
	}
	return count, nil
}

// GetRandomImage 随机获取图片
func GetRandomImage() (*model.Image, error) {
	var image model.Image
	query := db.Preload("MainTag").Preload("Tags").Order("RANDOM()")

	if err := query.First(&image).Error; err != nil {
		return nil, errors.WithStack(errs.ImageNotFound)
	}

	return &image, nil
}
