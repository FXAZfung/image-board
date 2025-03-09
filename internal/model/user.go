package model

import (
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/pkg/utils"
	"github.com/pkg/errors"
)

const (
	GENERAL = iota
	GUEST   // only one exists
	ADMIN   // only one exists
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`                      // unique key
	Username string `json:"username" gorm:"unique" binding:"required"` // username
	PwdHash  string `json:"-"`
	Role     int    `json:"role"` // user's role
	Disabled bool   `json:"disabled"`
	// Determine permissions by bit
	Permission int32 `json:"permission"` // password hash
}

// ValidatePwdStaticHash 验证密码是否正确
func (u *User) ValidatePwdStaticHash(password string) error {
	err := utils.ComparePassword(u.PwdHash, password)
	if err != nil {
		return errors.WithStack(errs.ErrUsernameOrPassword)
	}
	return nil
}

func (u *User) SetPassword(pwd string) *User {
	u.PwdHash, _ = utils.EncryptPassword(pwd)
	return u
}

func (u *User) IsGuest() bool {
	return u.Role == GUEST
}

func (u *User) IsAdmin() bool {
	return u.Role == ADMIN
}
