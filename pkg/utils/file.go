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

// ScanDir 扫描目录下的所有文件
func ScanDir(path string) []string {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Errorf("scan dir error: %s", err)
	}
	return files
}

// RemoveFile 删除文件
func RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		log.Errorf("remove file error: %s", err)
	}
	return err
}

// RemoveAll 删除文件夹
func RemoveAll(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		log.Errorf("remove all error: %s", err)
	}
	return err
}
