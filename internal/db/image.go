package db

import (
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/pkg/errors"
)

// GetImageByFilename 根据文件名获取图片
func GetImageByFilename(filename string) (*model.Image, error) {
	image := model.Image{FileName: filename}
	if err := db.Where(image).First(&image).Error; err != nil {
		return nil, errors.WithStack(errs.ImageNotFound)
	}
	return &image, nil
}

// GetImagesByPage GetImages 分页获取图片
func GetImagesByPage(page, pageSize int) ([]*model.Image, int64, error) {
	var images []*model.Image
	var count int64
	if err := db.Model(&model.Image{}).Count(&count).Error; err != nil {
		return nil, 0, errors.WithStack(errs.ErrImageCount)
	}
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&images).Error; err != nil {
		return nil, 0, errors.WithStack(errs.ErrImageList)
	}
	return images, count, nil
}

// GetImageByShortLink 根据短链获取图片
func GetImageByShortLink(shortLink string) (*model.Image, error) {
	image := model.Image{ShortLink: shortLink}
	if err := db.Where(image).First(&image).Error; err != nil {
		return nil, errors.WithStack(errs.ImageNotFound)
	}
	return &image, nil
}

// CreateImage 存储图片
func CreateImage(image *model.Image) error {
	if err := db.Save(image).Error; err != nil {
		return errors.WithStack(errs.ImageSaveError)
	}
	return nil
}

// GetImageCount 获取图片总数
func GetImageCount() (int64, error) {
	var count int64
	if err := db.Model(&model.Image{}).Count(&count).Error; err != nil {
		return 0, errors.WithStack(errs.ErrImageCount)
	}
	return count, nil
}

// GetRandomImage 随机获取一个图片
func GetRandomImage(category string) (*model.Image, error) {
	image := model.Image{}
	if category != "" {
		image.Category = category
	}
	//TODO 假如使用mysql 需要用 RAND()函数 使用 sqlite3 postgresql 需要使用RANDOM()函数
	if err := db.Where(image).Order("RANDOM()").First(&image).Error; err != nil {
		return nil, errors.WithStack(errs.ImageNotFound)
	}
	return &image, nil
}
