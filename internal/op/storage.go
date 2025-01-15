package op

import "github.com/FXAZfung/image-board/pkg/utils"

func GetStorageUsage(path string) (int64, error) {
	return utils.GetStorageUsage(path)
}
