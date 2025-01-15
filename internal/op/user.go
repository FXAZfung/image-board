package op

import (
	"github.com/FXAZfung/go-cache"
	"github.com/FXAZfung/image-board/internal/db"
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/pkg/singleflight"
	"time"
)

var userCache = cache.NewMemCache(cache.WithShards[*model.User](2))
var userG singleflight.Group[*model.User]
var guestUser *model.User
var adminUser *model.User

func GetAdmin() (*model.User, error) {
	if adminUser == nil {
		user, err := db.GetUserByRole(model.ADMIN)
		if err != nil {
			return nil, err
		}
		adminUser = user
	}
	return adminUser, nil
}

func GetGuest() (*model.User, error) {
	if guestUser == nil {
		user, err := db.GetUserByRole(model.GUEST)
		if err != nil {
			return nil, err
		}
		guestUser = user
	}
	return guestUser, nil
}

// GetUsers 分页获取用户
func GetUsers(page, perPage int) ([]model.User, int64, error) {
	users, total, err := db.GetUsers(page, perPage)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func GetUserByName(username string) (*model.User, error) {
	if username == "" {
		return nil, errs.ErrEmptyUsername
	}
	if user, ok := userCache.Get(username); ok {
		return user, nil
	}
	user, err, _ := userG.Do(username, func() (*model.User, error) {
		_user, err := db.GetUserByName(username)
		if err != nil {
			return nil, err
		}
		userCache.Set(username, _user, cache.WithEx[*model.User](time.Hour))
		return _user, nil
	})
	return user, err
}

func UpdateUser(u *model.User) error {
	old, err := db.GetUserById(u.ID)
	if err != nil {
		return err
	}
	userCache.Del(old.Username)
	return db.UpdateUser(u)
}

func CreateUser(u *model.User) error {
	userCache.Del(u.Username)
	return db.CreateUser(u)
}

// GetUserCount 获取用户总数
func GetUserCount() (int64, error) {
	count, err := db.GetUserCount()
	if err != nil {
		return 0, err
	}
	return count, nil
}
