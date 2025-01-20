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

func GetImagesByPage(pageIndex int, pageSize int) ([]*model.Image, int64, error) {
	images, count, err := db.GetImagesByPage(pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return images, count, nil
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

// GetImageCount 获取总的图片数量
func GetImageCount() (int64, error) {
	count, err := db.GetImageCount()
	if err != nil {
		return 0, err
	}
	return count, nil
}

//TODO: 重新从磁盘加载图片同时更新数据库
//// ReloadImages 重新从磁盘加载图片同时更新数据库
//func ReloadImages() error {
//	// 扫描磁盘
//	images := utils.ScanDir(config.Conf.DataImage.Dir)
//	// 如果数据库中没有的图片则添加
//
//	// 如果数据库中有但磁盘中没有的图片则删除
//
//	// 如果数据库中有且磁盘中也有的图片则更新
//
//	return nil
//}
