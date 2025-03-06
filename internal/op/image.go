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

var imageCache = cache.NewMemCache(cache.WithShards[*model.Image](4))
var imageListCache = cache.NewMemCache(cache.WithShards[interface{}](2))
var imageG singleflight.Group[*model.Image]
var imageListG singleflight.Group[interface{}]

// ImageCacheUpdate 清除图片缓存
func ImageCacheUpdate() {
	imageCache.Clear()
	imageListCache.Clear()
}

var imageCacheF = func(image *model.Image) {
	imageCache.Set(image.FileName, image, cache.WithEx[*model.Image](time.Minute*10))
	imageCache.Set(image.Hash, image, cache.WithEx[*model.Image](time.Minute*10))
	imageCache.Set(strconv.Itoa(int(image.ID)), image, cache.WithEx[*model.Image](time.Minute*10))
}

// GetImageByID 根据ID获取图片
func GetImageByID(id uint) (*model.Image, error) {
	key := strconv.Itoa(int(id))
	if image, ok := imageCache.Get(key); ok {
		return image, nil
	}

	image, err, _ := imageG.Do(key, func() (*model.Image, error) {
		_image, err := db.GetImageByID(id)
		if err != nil {
			return nil, err
		}
		imageCacheF(_image)
		return _image, nil
	})
	return image, err
}

// GetImageByFileName 根据文件名获取图片
func GetImageByFileName(fileName string) (*model.Image, error) {
	if image, ok := imageCache.Get(fileName); ok {
		return image, nil
	}

	image, err, _ := imageG.Do(fileName, func() (*model.Image, error) {
		_image, err := db.GetImageByFilename(fileName)
		if err != nil {
			return nil, err
		}
		imageCacheF(_image)
		return _image, nil
	})
	return image, err
}

// GetImageByHash 根据hash获取图片
func GetImageByHash(hash string) (*model.Image, error) {
	if image, ok := imageCache.Get(hash); ok {
		return image, nil
	}

	image, err, _ := imageG.Do(hash, func() (*model.Image, error) {
		_image, err := db.GetImageByHash(hash)
		if err != nil {
			return nil, err
		}
		imageCacheF(_image)
		return _image, nil
	})
	return image, err
}

// GetImagesByPage 分页获取图片
func GetImagesByPage(page, pageSize int) ([]*model.Image, int64, error) {
	cacheKey := fmt.Sprintf("images_page_%d_%d", page, pageSize)
	if cached, ok := imageListCache.Get(cacheKey); ok {
		data := cached.(map[string]interface{})
		return data["images"].([]*model.Image), data["count"].(int64), nil
	}

	result, err, _ := imageListG.Do(cacheKey, func() (interface{}, error) {
		images, count, err := db.GetImagesByPage(page, pageSize)
		if err != nil {
			return nil, err
		}

		// Cache individual images
		for _, img := range images {
			imageCacheF(img)
		}

		// Cache the page result
		data := map[string]interface{}{
			"images": images,
			"count":  count,
		}
		imageListCache.Set(cacheKey, data, cache.WithEx[interface{}](time.Minute*2))
		return data, nil
	})

	if result != nil {
		data := result.(map[string]interface{})
		return data["images"].([]*model.Image), data["count"].(int64), nil
	}
	return nil, 0, err
}

// GetImagesByTag 获取拥有特定标签的所有图片
func GetImagesByTag(tagName string, page, pageSize int) ([]*model.Image, int64, error) {
	cacheKey := fmt.Sprintf("images_tag_%s_%d_%d", tagName, page, pageSize)
	if cached, ok := imageListCache.Get(cacheKey); ok {
		data := cached.(map[string]interface{})
		return data["images"].([]*model.Image), data["count"].(int64), nil
	}

	result, err, _ := imageListG.Do(cacheKey, func() (interface{}, error) {
		images, count, err := db.GetImagesByTag(tagName, page, pageSize)
		if err != nil {
			return nil, err
		}

		// Cache individual images
		for _, img := range images {
			imageCacheF(img)
		}

		// Cache the tag-based image list
		data := map[string]interface{}{
			"images": images,
			"count":  count,
		}
		imageListCache.Set(cacheKey, data, cache.WithEx[interface{}](time.Minute*2))
		return data, nil
	})

	if result != nil {
		data := result.(map[string]interface{})
		return data["images"].([]*model.Image), data["count"].(int64), nil
	}
	return nil, 0, err
}

// CreateImage 创建新图片
func CreateImage(image *model.Image) error {
	if err := db.CreateImage(image); err != nil {
		return err
	}
	imageCacheF(image)
	// 清除可能受影响的列表缓存
	imageListCache.Clear()
	return nil
}

// UpdateImage 更新图片信息
func UpdateImage(image *model.Image) error {
	if err := db.UpdateImage(image); err != nil {
		return err
	}
	// 更新缓存
	imageCacheF(image)
	return nil
}

// DeleteImage 删除图片
func DeleteImage(imageID uint) error {
	// 获取图片以便缓存失效
	image, err := GetImageByID(imageID)
	if err == nil {
		imageCache.Del(image.FileName)
		imageCache.Del(image.Hash)
		imageCache.Del(strconv.Itoa(int(image.ID)))
	}

	if err := db.DeleteImage(imageID); err != nil {
		return err
	}

	// 清除可能受影响的列表缓存
	imageListCache.Clear()
	return nil
}

// GetRandomImage 获取随机图片
func GetRandomImage() (*model.Image, error) {
	image, err := db.GetRandomImage()
	if err != nil {
		return nil, err
	}
	// 缓存随机获取的图片
	imageCacheF(image)
	return image, nil
}

// GetImageCount 获取总的图片数量
func GetImageCount() (int64, error) {
	cacheKey := "image_count"
	if cached, ok := imageListCache.Get(cacheKey); ok {
		return cached.(int64), nil
	}

	result, err, _ := imageListG.Do(cacheKey, func() (interface{}, error) {
		count, err := db.GetImageCount()
		if err != nil {
			return int64(0), err
		}
		imageListCache.Set(cacheKey, count, cache.WithEx[interface{}](time.Minute*5))
		return count, nil
	})

	if result != nil {
		return result.(int64), err
	}
	return 0, err
}
