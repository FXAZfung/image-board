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

// GetImagesByPage retrieves paginated images with their tags
func GetImagesByPage(page, perPage int) ([]*model.Image, int64, error) {
	var count int64
	if err := db.Model(&model.Image{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var images []*model.Image
	offset := (page - 1) * perPage

	// Add Preload("Tags") to load the associated tags
	if err := db.Preload("Tags").Offset(offset).Limit(perPage).Order("id desc").Find(&images).Error; err != nil {
		return nil, 0, err
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
func RemoveTagFromImage(imageID uint, tagID uint) (*model.Tag, error) {
	// 先获取标签信息，用户返回，然后执行事务
	var tag model.Tag
	if err := db.First(&tag, tagID).Error; err != nil {
		return nil, errors.WithStack(errs.ErrTagNotFound)
	}
	db.Transaction(func(tx *gorm.DB) error {
		// 检查图片是否存在
		var image model.Image
		if err := tx.First(&image, imageID).Error; err != nil {
			return errors.WithStack(errs.ImageNotFound)
		}

		// 检查标签是否存在
		var tag model.Tag
		if err := tx.First(&tag, tagID).Error; err != nil {
			return errors.WithStack(errs.ErrTagNotFound)
		}

		// 检查图片和标签的关联是否存在
		var imageTag model.ImageTag
		if err := tx.Where("image_id = ? AND tag_id = ?", imageID, tagID).First(&imageTag).Error; err != nil {
			return errors.WithStack(errs.ErrTagNotFound)
		}

		// 删除标签和图片的关联
		if err := tx.Delete(&imageTag).Error; err != nil {
			return errors.WithStack(err)
		}

		// 减少标签计数
		if err := tx.Model(&tag).Update("count", tag.Count-1).Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	return &tag, nil
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
	query := db.Preload("Tags").Order("RANDOM()")

	if err := query.First(&image).Error; err != nil {
		return nil, errors.WithStack(errs.ImageNotFound)
	}

	return &image, nil
}
