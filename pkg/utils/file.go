package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// IsExist 判断文件是否存在
func IsExist(path string) bool {
	_, err :=
		os.Stat(path)
	return err == nil || os.IsExist(err)
}

// CreateNestedFile create nested file
func CreateNestedFile(path string) (*os.File, error) {
	basePath := filepath.Dir(path)
	if err := CreateNestedDirectory(basePath); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// CreateNestedDirectory create nested directory
func CreateNestedDirectory(path string) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		log.Errorf("can't create folder, %s", err)
	}
	return err
}
