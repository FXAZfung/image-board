package data

import (
	"github.com/FXAZfung/image-board/cmd/flags"
	"github.com/FXAZfung/image-board/internal/db"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/pkg/random"
	"github.com/FXAZfung/image-board/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"os"
)

func initUser() {
	admin, err := op.GetAdmin()
	adminPassword := random.String(8)
	envpass := os.Getenv("IMAGE_BOARD_ADMIN_PASSWORD")
	if flags.Dev {
		adminPassword = "admin"
	} else if len(envpass) > 0 {
		adminPassword = envpass
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pwdHash, err := utils.EncryptPassword(adminPassword)
			if err != nil {
				utils.Log.Fatalf("[init user] Failed to encrypt admin password %v", err)
			}
			admin = &model.User{
				Username: "admin",
				PwdHash:  pwdHash,
				Role:     model.ADMIN,
			}
			if err := op.CreateUser(admin); err != nil {
				panic(err)
			} else {
				utils.Log.Infof("Successfully created the admin user and the initial password is: %s", adminPassword)
			}
		} else {
			utils.Log.Fatalf("[init user] Failed to get admin user: %v", err)
		}
	}
	guest, err := op.GetGuest()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pwdHash, err := utils.EncryptPassword("guest")
			if err != nil {
				utils.Log.Fatalf("[init user] Failed to encrypt guest password %v", err)
			}
			guest = &model.User{
				Username:   "guest",
				PwdHash:    pwdHash,
				Role:       model.GUEST,
				Permission: 0,
				Disabled:   true,
			}
			if err := db.CreateUser(guest); err != nil {
				utils.Log.Fatalf("[init user] Failed to create guest user: %v", err)
			}
		} else {
			utils.Log.Fatalf("[init user] Failed to get guest user: %v", err)
		}
	}
}
