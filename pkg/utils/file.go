package utils

import (
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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

// GetStorageUsage 获取存储使用情况
func GetStorageUsage(path string) (int64, error) {

	// 检查路径是否存在
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	// 如果是文件，直接返回文件大小
	if !info.IsDir() {
		return info.Size(), nil
	}

	// 如果是目录，递归计算目录大小
	var totalSize int64
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	return totalSize, err
}

// IsImage 判断文件是否为图片
func IsImage(file *multipart.FileHeader) bool {
	// 检查 Content-Type
	contentType := file.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif",
		"image/webp", "image/svg+xml", "image/tiff", "image/bmp":
		// 继续检查文件扩展名作为二次验证
	default:
		return false
	}

	// 检查文件扩展名
	filename := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".svg": true, ".tiff": true, ".bmp": true,
	}

	// 同时验证 Content-Type 和扩展名
	return validExtensions[ext]
}
