package op

import (
	"fmt"
	"github.com/FXAZfung/go-cache"
	"github.com/FXAZfung/image-board/internal/db"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/pkg/singleflight"
	"strconv"
	"time"
)

var tagCache = cache.NewMemCache(cache.WithShards[*model.Tag](4))
var tagListCache = cache.NewMemCache(cache.WithShards[interface{}](2))
var tagG singleflight.Group[*model.Tag]
var tagListG singleflight.Group[interface{}]

// TagCacheUpdate clears all tag caches
func TagCacheUpdate() {
	tagCache.Clear()
	tagListCache.Clear()
}

// cacheTag helper to store tag in cache
var cacheTag = func(tag *model.Tag) {
	tagCache.Set(tag.Name, tag, cache.WithEx[*model.Tag](time.Minute*10))
	tagCache.Set(strconv.Itoa(int(tag.ID)), tag, cache.WithEx[*model.Tag](time.Minute*10))
}

// GetTagByID retrieves a tag by its ID with caching
func GetTagByID(id uint) (*model.Tag, error) {
	key := strconv.Itoa(int(id))
	if tag, ok := tagCache.Get(key); ok {
		return tag, nil
	}

	tag, err, _ := tagG.Do(key, func() (*model.Tag, error) {
		tag, err := db.GetTagByID(id)
		if err != nil {
			return nil, err
		}
		cacheTag(tag)
		return tag, nil
	})
	return tag, err
}

// GetTagByName retrieves a tag by its name with caching
func GetTagByName(name string) (*model.Tag, error) {
	if tag, ok := tagCache.Get(name); ok {
		return tag, nil
	}

	tag, err, _ := tagG.Do(name, func() (*model.Tag, error) {
		tag, err := db.GetTagByName(name)
		if err != nil {
			return nil, err
		}
		cacheTag(tag)
		return tag, nil
	})
	return tag, err
}

// GetOrCreateTag gets a tag by name or creates it if it doesn't exist
func GetOrCreateTag(name string) (*model.Tag, bool, error) {
	// Try cache first
	if tag, ok := tagCache.Get(name); ok {
		return tag, false, nil
	}

	// Use singleflight to prevent duplicate creations
	result, err, _ := tagG.Do("create_"+name, func() (*model.Tag, error) {
		tag, created, err := db.GetOrCreateTag(name)
		if err != nil {
			return nil, err
		}
		cacheTag(tag)
		// If created, invalidate list cache
		if created {
			tagListCache.Clear()
		}
		return tag, nil
	})

	// Need to get created status from database again
	_, created, _ := db.GetOrCreateTag(name)
	return result, created, err
}

// ListTags retrieves all tags with pagination and caching
func ListTags(page, pageSize int) ([]*model.Tag, int64, error) {
	cacheKey := fmt.Sprintf("tags_page_%d_%d", page, pageSize)
	if cached, ok := tagListCache.Get(cacheKey); ok {
		data := cached.(map[string]interface{})
		return data["tags"].([]*model.Tag), data["count"].(int64), nil
	}

	result, err, _ := tagListG.Do(cacheKey, func() (interface{}, error) {
		tags, count, err := db.ListTags(page, pageSize)
		if err != nil {
			return nil, err
		}

		// Cache individual tags
		for _, tag := range tags {
			cacheTag(tag)
		}

		// Cache the page result
		data := map[string]interface{}{
			"tags":  tags,
			"count": count,
		}
		tagListCache.Set(cacheKey, data, cache.WithEx[interface{}](time.Minute*5))
		return data, nil
	})

	if result != nil {
		data := result.(map[string]interface{})
		return data["tags"].([]*model.Tag), data["count"].(int64), nil
	}
	return nil, 0, err
}

// CreateTag creates a new tag with cache handling
func CreateTag(tag *model.Tag) error {
	if err := db.CreateTag(tag); err != nil {
		return err
	}
	cacheTag(tag)
	// Clear list cache since we added a new tag
	tagListCache.Clear()
	return nil
}

// UpdateTag updates a tag's information and updates cache
func UpdateTag(tag *model.Tag) error {
	if err := db.UpdateTag(tag); err != nil {
		return err
	}
	// Update cache
	cacheTag(tag)
	return nil
}

// DeleteTag deletes a tag and removes it from cache
func DeleteTag(tagID uint) error {
	// Get tag to invalidate cache
	tag, err := GetTagByID(tagID)
	if err == nil {
		tagCache.Del(tag.Name)
		tagCache.Del(strconv.Itoa(int(tag.ID)))
	}

	if err := db.DeleteTag(tagID); err != nil {
		return err
	}

	// Clear list cache
	tagListCache.Clear()
	return nil
}

// AddTagsToImage adds tags to an image with cache updates
func AddTagsToImage(imageID uint, tags []string) error {
	if err := db.AddTagsToImage(imageID, tags); err != nil {
		return err
	}

	// Invalidate tag caches since counts may have changed
	tagListCache.Clear()

	// Invalidate image cache for this image
	ImageCacheUpdate()

	return nil
}

// RemoveTagFromImage removes a tag from an image
func RemoveTagFromImage(imageID uint, tagID uint) error {
	if err := db.RemoveTagFromImage(imageID, tagID); err != nil {
		return err
	}

	// Get tag to update cache
	tag, err := GetTagByID(tagID)
	if err == nil {
		// Refresh tag in cache since count changed
		cacheTag(tag)
	} else {
		// If error, just invalidate tag cache
		tagCache.Del(strconv.Itoa(int(tagID)))
	}

	// Invalidate lists that may contain this tag
	tagListCache.Clear()

	// Invalidate image cache for this image
	key := strconv.Itoa(int(imageID))
	imageCache.Del(key)

	return nil
}

// GetMostPopularTags retrieves the most used tags with caching
func GetMostPopularTags(limit int) ([]*model.Tag, error) {
	cacheKey := fmt.Sprintf("popular_tags_%d", limit)
	if cached, ok := tagListCache.Get(cacheKey); ok {
		return cached.([]*model.Tag), nil
	}

	result, err, _ := tagListG.Do(cacheKey, func() (interface{}, error) {
		tags, err := db.GetMostPopularTags(limit)
		if err != nil {
			return nil, err
		}

		// Cache individual tags
		for _, tag := range tags {
			cacheTag(tag)
		}

		tagListCache.Set(cacheKey, tags, cache.WithEx[interface{}](time.Minute*5))
		return tags, nil
	})

	if result != nil {
		return result.([]*model.Tag), nil
	}
	return nil, err
}

// GetTagsForImage retrieves all tags for a specific image with caching
func GetTagsForImage(imageID uint) ([]*model.Tag, error) {
	cacheKey := fmt.Sprintf("image_tags_%d", imageID)
	if cached, ok := tagListCache.Get(cacheKey); ok {
		return cached.([]*model.Tag), nil
	}

	result, err, _ := tagListG.Do(cacheKey, func() (interface{}, error) {
		tags, err := db.GetTagsForImage(imageID)
		if err != nil {
			return nil, err
		}

		// Cache individual tags
		for _, tag := range tags {
			cacheTag(tag)
		}

		tagListCache.Set(cacheKey, tags, cache.WithEx[interface{}](time.Minute*5))
		return tags, nil
	})

	if result != nil {
		return result.([]*model.Tag), nil
	}
	return nil, err
}

// SearchTagsByPrefix searches tags that start with the specified prefix
func SearchTagsByPrefix(prefix string, limit int) ([]*model.Tag, error) {
	cacheKey := fmt.Sprintf("tag_prefix_%s_%d", prefix, limit)
	if cached, ok := tagListCache.Get(cacheKey); ok {
		return cached.([]*model.Tag), nil
	}

	result, err, _ := tagListG.Do(cacheKey, func() (interface{}, error) {
		tags, err := db.SearchTagsByPrefix(prefix, limit)
		if err != nil {
			return nil, err
		}

		// Cache individual tags
		for _, tag := range tags {
			cacheTag(tag)
		}

		tagListCache.Set(cacheKey, tags, cache.WithEx[interface{}](time.Minute*1))
		return tags, nil
	})

	if result != nil {
		return result.([]*model.Tag), nil
	}
	return nil, err
}
