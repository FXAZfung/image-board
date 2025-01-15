package db

import (
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/pkg/errors"
)

func GetUserByRole(role int) (*model.User, error) {
	user := model.User{Role: role}
	if err := db.Where(user).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByName 根据用户名获取用户
func GetUserByName(username string) (*model.User, error) {
	user := model.User{Username: username}
	if err := db.Where(user).First(&user).Error; err != nil {
		return nil, errors.Wrapf(err, "failed find user")
	}
	return &user, nil
}

// GetUserById 根据用户ID获取用户
func GetUserById(id uint) (*model.User, error) {
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed get old user")
	}
	return &user, nil
}

func UpdateUser(u *model.User) error {
	return errors.WithStack(db.Save(u).Error)
}

func CreateUser(u *model.User) error {
	return errors.WithStack(db.Create(u).Error)
}

// GetUserCount 获取用户数量
func GetUserCount() (int64, error) {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, errors.Wrapf(err, "failed get user count")
	}
	return count, nil
}

// GetUsers 分页获取用户
func GetUsers(pageIndex, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get user count")
	}
	if err := db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get users")
	}
	return users, count, nil
}
