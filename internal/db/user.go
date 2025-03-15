package db

import (
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GetUserById retrieves a user by their ID
func GetUserById(id uint) (*model.User, error) {
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(errs.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return &user, nil
}

// GetUserByName retrieves a user by their username
func GetUserByName(username string) (*model.User, error) {
	var user model.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(errs.ErrUserNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return &user, nil
}

// GetUserByRole retrieves a user by their role (admin, guest, etc.)
func GetUserByRole(role int) (*model.User, error) {
	var user model.User
	if err := db.Where("role = ?", role).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(gorm.ErrRecordNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return &user, nil
}

// GetUsers retrieves users with pagination
func GetUsers(page, perPage int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	if err := db.Order("id ASC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&users).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return users, total, nil
}

// CreateUser creates a new user in the database
func CreateUser(user *model.User) error {
	// Check if user with same username already exists
	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return errors.WithStack(err)
	}

	if count > 0 {
		return errors.WithStack(errs.ErrUserExist)
	}

	return db.Create(user).Error
}

// UpdateUser updates an existing user's information
func UpdateUser(user *model.User) error {
	return db.Save(user).Error
}

// DeleteUser deletes a user by their ID
func DeleteUser(id uint) error {
	return db.Delete(&model.User{}, id).Error
}

// GetUserCount returns the total number of users
func GetUserCount() (int64, error) {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, errors.WithStack(err)
	}
	return count, nil
}
