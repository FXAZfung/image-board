package op

import (
	"github.com/FXAZfung/go-cache"
	"github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/db"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/pkg/singleflight"
	"os"
	"path"
	"time"
)

var categoryCache = cache.NewMemCache(cache.WithShards[*model.Category](4))
var categoryCacheG singleflight.Group[*model.Category]
var categoryCacheF = func(item *model.Category) {
	categoryCache.Set(item.Name, item, cache.WithEx[*model.Category](time.Hour))
}

var categoryGroupCache = cache.NewMemCache(cache.WithShards[[]*model.Category](4))
var categoryGroupG singleflight.Group[[]*model.Category]
var categoryGroupCacheF = func(key string, item []*model.Category) {
	categoryGroupCache.Set(key, item, cache.WithEx[[]*model.Category](time.Hour))
}

func categoryCacheUpdate() {
	categoryCache.Clear()
	categoryGroupCache.Clear()
}

func GetCategoryByName(name string) (*model.Category, error) {
	if item, ok := categoryCache.Get(name); ok {
		return item, nil
	}

	item, err, _ := categoryCacheG.Do(name, func() (*model.Category, error) {
		_item, err := db.GetCategoryByName(name)
		if err != nil {
			return nil, err
		}
		categoryCacheF(_item)
		return _item, nil
	})
	return item, err
}

func GetCategories() ([]*model.Category, error) {
	if items, ok := categoryGroupCache.Get("ALL_CATEGORIES"); ok {
		return items, nil
	}
	items, err, _ := categoryGroupG.Do("ALL_CATEGORIES", func() ([]*model.Category, error) {
		_items, err := db.GetCategories()
		if err != nil {
			return nil, err
		}
		categoryGroupCacheF("ALL_CATEGORIES", _items)
		return _items, nil
	})
	return items, err
}

func SaveCategory(item *model.Category) (err error) {
	// hook
	if _, err := HandleCategoryHook(item); err != nil {
		return err
	}
	// 创建分类文件夹
	dir := path.Join(config.Conf.DataImage.Dir, item.Name)
	err = os.MkdirAll(dir, 0755)
	err = db.CreateCategory(item)
	if err != nil {
		return err
	}
	categoryCacheUpdate()
	return nil
}

// GetCategoryCount 获取分类总数
func GetCategoryCount() (int64, error) {
	count, err := db.GetCategoryCount()
	if err != nil {
		return 0, err
	}
	return count, nil
}
