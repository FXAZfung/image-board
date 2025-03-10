package db

import (
	"fmt"
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateTag creates a new tag in the database
func CreateTag(tag *model.Tag) error {
	return db.Create(tag).Error
}

// GetTagByID retrieves a tag by its ID
func GetTagByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := db.First(&tag, id).Error; err != nil {
		return nil, errors.WithStack(errs.ErrTagNotFound)
	}
	return &tag, nil
}

// GetTagByName retrieves a tag by its name
func GetTagByName(name string) (*model.Tag, error) {
	var tag model.Tag
	if err := db.Where("name = ?", name).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(errs.ErrTagNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return &tag, nil
}

// GetOrCreateTag gets an existing tag or creates a new one if it doesn't exist
func GetOrCreateTag(name string) (*model.Tag, bool, error) {
	var tag model.Tag
	created := false

	// Try to find the tag first
	err := db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new tag if not found
			tag = model.Tag{
				Name:  name,
				Count: 0,
			}
			if err := db.Create(&tag).Error; err != nil {
				return nil, false, errors.WithStack(err)
			}
			created = true
		} else {
			return nil, false, errors.WithStack(err)
		}
	}

	return &tag, created, nil
}

// ListTags retrieves all tags with pagination
func ListTags(page, pageSize int) ([]*model.Tag, int64, error) {
	var tags []*model.Tag
	var count int64

	if err := db.Model(&model.Tag{}).Count(&count).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	if err := db.Order("count DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tags).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return tags, count, nil
}

// UpdateTag updates a tag's information
func UpdateTag(tag *model.Tag) error {
	return db.Save(tag).Error
}

// DeleteTag deletes a tag and removes its associations with images
func DeleteTag(tagID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// First delete associations
		if err := tx.Where("tag_id = ?", tagID).Delete(&model.ImageTag{}).Error; err != nil {
			return errors.WithStack(err)
		}

		// Then delete the tag itself
		if err := tx.Delete(&model.Tag{}, tagID).Error; err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}

// AddTagToImage adds a single tag to an image
func AddTagToImage(imageID uint, tagName string) error {
	return AddTagsToImage(imageID, []string{tagName})
}

// AddTagsToImage adds tags to an image and handles the tag creation if needed
func AddTagsToImage(imageID uint, tagNames []string) error {
	if len(tagNames) == 0 {
		return nil // Nothing to do if no tags provided
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// First verify the image exists
		var image model.Image
		if err := tx.First(&image, imageID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.WithStack(errs.ImageNotFound)
			}
			return errors.WithStack(err)
		}

		// Process each tag
		for _, name := range tagNames {
			if name == "" {
				continue // Skip empty tag names
			}

			// Get or create the tag
			var tag model.Tag
			result := tx.Where("name = ?", name).FirstOrCreate(&tag, model.Tag{
				Name:  name,
				Count: 0, // Will be incremented below
			})
			if result.Error != nil {
				return errors.WithStack(result.Error)
			}

			// Check if association already exists to avoid duplicates
			var count int64
			if err := tx.Model(&model.ImageTag{}).
				Where("image_id = ? AND tag_id = ?", imageID, tag.ID).
				Count(&count).Error; err != nil {
				return errors.WithStack(err)
			}

			// Skip if already associated
			if count > 0 {
				continue
			}

			// Create the association
			imageTag := model.ImageTag{
				ImageID: imageID,
				TagID:   tag.ID,
			}
			if err := tx.Create(&imageTag).Error; err != nil {
				return errors.WithStack(err)
			}

			// Increment tag count
			if err := tx.Model(&tag).Update("count", gorm.Expr("count + 1")).Error; err != nil {
				return errors.WithStack(err)
			}

			log.Debugf("Added tag %s (ID: %d) to image %d", name, tag.ID, imageID)
		}

		return nil
	})
}

// GetMostPopularTags retrieves the most used tags
func GetMostPopularTags(limit int) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := db.Order("count DESC").Limit(limit).Find(&tags).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return tags, nil
}

// GetTagsForImage retrieves all tags for a specific image
func GetTagsForImage(imageID uint) ([]*model.Tag, error) {
	var tags []*model.Tag

	err := db.Table("im_tags").
		Joins("JOIN im_image_tags ON im_image_tags.tag_id = im_tags.id").
		Where("im_image_tags.image_id = ?", imageID).
		Find(&tags).Error

	if err != nil {
		log.WithFields(log.Fields{
			"image_id": imageID,
			"error":    err,
		}).Error("Failed to retrieve tags for image")
		return nil, errors.WithStack(err)
	}

	return tags, nil
}

// SearchTagsByPrefix searches tags that start with the specified prefix
func SearchTagsByPrefix(prefix string, limit int) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := db.Where("name LIKE ?", fmt.Sprintf("%s%%", prefix)).Limit(limit).Find(&tags).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return tags, nil
}
