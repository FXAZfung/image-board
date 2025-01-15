package model

import (
	"github.com/FXAZfung/image-board/internal/errs"
	"github.com/FXAZfung/image-board/pkg/utils"
	"github.com/pkg/errors"
)

const (
	GENERAL = iota
	GUEST   // only one exists
	ADMIN
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`                      // unique key
	Username string `json:"username" gorm:"unique" binding:"required"` // username
	PwdHash  string `json:"-"`
	Role     int    `json:"role"` // user's role
	Disabled bool   `json:"disabled"`
	// Determine permissions by bit
	//   0: can see hidden files
	//   1: can access without password
	//   2: can add offline download tasks
	//   3: can mkdir and upload
	//   4: can rename
	//   5: can move
	//   6: can copy
	//   7: can remove
	//   8: webdav read
	//   9: webdav write
	Permission int32 `json:"permission"` // password hash
}

// ValidatePwdStaticHash 验证密码是否正确
func (u *User) ValidatePwdStaticHash(password string) error {
	//reqPassword, err := utils.EncryptPassword(password)
	//if err != nil {
	//	return errors.WithStack(errs.Internal)
	//}
	err := utils.ComparePassword(u.PwdHash, password)
	if err != nil {
		return errors.WithStack(errs.ErrUserPassword)
	}
	return nil
}

func (u *User) SetPassword(pwd string) *User {
	u.PwdHash, _ = utils.EncryptPassword(pwd)
	return u
}
