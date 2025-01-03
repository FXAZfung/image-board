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
