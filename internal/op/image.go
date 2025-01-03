package op

import (
	"github.com/FXAZfung/go-cache"
	"github.com/FXAZfung/image-board/internal/db"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/pkg/singleflight"
	"os"
	"time"
)

var imageCache = cache.NewMemCache(cache.WithShards[*model.Image](4))
var imageG singleflight.Group[*model.Image]

var imageCacheF = func(image *model.Image) {
	imageCache.Set(image.FileName, image, cache.WithEx[*model.Image](time.Minute*10))
}

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

func GetImagesByPage(page int, pageSize int) ([]*model.Image, error) {
	images, err := db.GetImagesByPage(page, pageSize)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func GetImageByShortLink(shortLink string) (*model.Image, error) {
	image, err := db.GetImageByShortLink(shortLink)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func CreateImage(image *model.Image, data []byte, path string) error {
	// 将图片保存到文件
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	err = db.CreateImage(image)
	if err != nil {
		return err
	}
	return nil
}

func GetRandomImage(category string) (*model.Image, error) {
	image, err := db.GetRandomImage(category)
	if err != nil {
		return nil, err
	}
	return image, nil
}
