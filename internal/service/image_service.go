package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/internal/model/request"
	"github.com/FXAZfung/image-board/internal/model/response"
	"github.com/FXAZfung/image-board/internal/op"
	"github.com/FXAZfung/image-board/pkg/utils"
	"github.com/disintegration/imaging"
	"golang.org/x/sync/errgroup"
)

const (
	maxConcurrentProcessing = 3 // 根据CPU核心数调整
)

var (
	processingSem = make(chan struct{}, maxConcurrentProcessing)
)

type ImageService struct {
	baseDir        string
	thumbnailWidth int
	quality        int
	allowedExts    []string
}

func NewImageService() *ImageService {
	if config.Conf.DataImage.Dir == "" {
		log.Fatal("Image storage directory not configured")
	}

	return &ImageService{
		baseDir:        config.Conf.DataImage.Dir,
		thumbnailWidth: 300,
		quality:        90,
		allowedExts:    []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	}
}

// UploadImage 入口函数
func UploadImage(file *multipart.FileHeader, user *model.User) (*model.Image, error) {
	startTime := time.Now()
	logFields := log.Fields{
		"user_id":   user.ID,
		"file_name": file.Filename,
		"file_size": file.Size,
		"operation": "upload",
	}

	defer func() {
		log.WithFields(logFields).
			WithField("duration", time.Since(startTime)).
			Info("Upload processing completed")
	}()

	service := NewImageService()
	image, err := service.processUpload(file, user, logFields)
	if err != nil {
		log.WithFields(logFields).Errorf("Upload failed: %v", err)
		return nil, err
	}

	return image, nil
}

type uploadContext struct {
	file          *multipart.FileHeader
	user          *model.User
	fileData      []byte
	hash          string
	fileExt       string
	filePath      string
	thumbnailPath string
	webpPath      string
	modImage      *model.Image
	logFields     log.Fields
}

func (s *ImageService) processUpload(file *multipart.FileHeader, user *model.User, logFields log.Fields) (*model.Image, error) {
	ctx := &uploadContext{
		file:      file,
		user:      user,
		logFields: logFields,
	}

	steps := []func() error{
		ctx.validateInput,
		ctx.readAndHashFile,
		ctx.checkDuplicate,
		ctx.validateExtension,
		ctx.generateFilePaths,
		ctx.createStorageDirs,
		ctx.processImageData,
		ctx.createImageModel,
		ctx.saveToDatabase,
	}

	for _, step := range steps {
		if err := step(); err != nil {
			return nil, s.wrapError("upload processing failed", err)
		}
	}

	return ctx.modImage, nil
}

func (ctx *uploadContext) validateInput() error {
	if ctx.file == nil {
		return errors.New("no file provided")
	}
	if !utils.IsImage(ctx.file) {
		return errors.New("invalid file type")
	}
	return nil
}

func (ctx *uploadContext) readAndHashFile() error {
	hash := sha256.New()
	f, err := ctx.file.Open()
	if err != nil {
		return fmt.Errorf("file open failed: %w", err)
	}
	defer f.Close()

	// 流式读取同时计算哈希
	tee := io.TeeReader(f, hash)
	data, err := io.ReadAll(tee)
	if err != nil {
		return fmt.Errorf("file read failed: %w", err)
	}

	ctx.fileData = data
	ctx.hash = hex.EncodeToString(hash.Sum(nil))
	return nil
}

func (ctx *uploadContext) checkDuplicate() error {
	if existing, err := op.GetImageByHash(ctx.hash); err == nil {
		log.WithFields(ctx.logFields).Info("Duplicate image found")
		ctx.modImage = existing
		return errors.New("duplicate image") // 特殊错误类型触发提前返回
	}
	return nil
}

func (s *ImageService) wrapError(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, errors.Cause(err))
}

func (ctx *uploadContext) validateExtension() error {
	ctx.fileExt = strings.ToLower(filepath.Ext(ctx.file.Filename))
	for _, ext := range NewImageService().allowedExts {
		if ctx.fileExt == ext {
			return nil
		}
	}
	return fmt.Errorf("invalid file extension: %s", ctx.fileExt)
}

func (ctx *uploadContext) generateFilePaths() error {
	now := time.Now()
	datePath := fmt.Sprintf("%d/%02d", now.Year(), now.Month())
	baseDir := filepath.Join(NewImageService().baseDir, datePath)

	ctx.filePath = filepath.Join(baseDir, ctx.hash+ctx.fileExt)
	ctx.thumbnailPath = filepath.Join(baseDir, "thumbnails", ctx.hash+ctx.fileExt)
	ctx.webpPath = GetWebPPath(filepath.Join(baseDir, "webp", ctx.hash+ctx.fileExt))
	return nil
}

func (ctx *uploadContext) createStorageDirs() error {
	dirs := []string{
		filepath.Dir(ctx.filePath),
		filepath.Dir(ctx.thumbnailPath),
		filepath.Dir(ctx.webpPath),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("directory creation failed: %w", err)
		}
	}
	return nil
}

func (ctx *uploadContext) processImageData() error {
	g, _ := errgroup.WithContext(context.Background())

	// 保存原始文件
	g.Go(func() error {
		return safeWriteFile(ctx.filePath, ctx.fileData)
	})

	// 生成缩略图
	g.Go(func() error {
		processingSem <- struct{}{}
		defer func() { <-processingSem }()
		return NewImageService().createThumbnail(ctx.fileData, ctx.thumbnailPath)
	})

	// 生成WebP
	g.Go(func() error {
		processingSem <- struct{}{}
		defer func() { <-processingSem }()
		return NewImageService().convertToWebP(ctx.fileData, ctx.webpPath)
	})

	if err := g.Wait(); err != nil {
		ctx.cleanupFiles()
		return fmt.Errorf("file processing failed: %w", err)
	}
	return nil
}

func safeWriteFile(path string, data []byte) error {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func (ctx *uploadContext) createImageModel() error {
	img, _, err := image.Decode(bytes.NewReader(ctx.fileData))
	if err != nil {
		return fmt.Errorf("image decode failed: %w", err)
	}

	ctx.modImage = &model.Image{
		FileName:      ctx.hash + ctx.fileExt,
		OriginalName:  ctx.file.Filename,
		Hash:          ctx.hash,
		ContentType:   http.DetectContentType(ctx.fileData),
		Size:          ctx.file.Size,
		Path:          ctx.filePath,
		ThumbnailPath: ctx.thumbnailPath,
		WebpPath:      ctx.webpPath,
		Width:         img.Bounds().Dx(),
		Height:        img.Bounds().Dy(),
		UserID:        ctx.user.ID,
		IsPublic:      true,
	}
	return nil
}

func (ctx *uploadContext) saveToDatabase() error {
	if err := op.CreateImage(ctx.modImage); err != nil {
		ctx.cleanupFiles()
		return fmt.Errorf("database save failed: %w", err)
	}
	return nil
}

func (ctx *uploadContext) cleanupFiles() {
	files := []string{ctx.filePath, ctx.thumbnailPath, ctx.webpPath}
	var wg sync.WaitGroup

	for _, path := range files {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				log.WithFields(ctx.logFields).Warnf("Cleanup failed for %s: %v", p, err)
			}
		}(path)
	}
	wg.Wait()
}

// 其他方法保持类似结构，以下是修改后的关键函数：

func (s *ImageService) createThumbnail(imgData []byte, path string) error {
	src, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return fmt.Errorf("thumbnail decode failed: %w", err)
	}

	thumbnail := imaging.Resize(src, s.thumbnailWidth, 0, imaging.Lanczos)
	if err := imaging.Save(thumbnail, path, imaging.JPEGQuality(s.quality)); err != nil {
		return fmt.Errorf("thumbnail save failed: %w", err)
	}
	return nil
}

func (s *ImageService) convertToWebP(imgData []byte, path string) error {
	src, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return fmt.Errorf("webp decode failed: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("webp create failed: %w", err)
	}
	defer f.Close()

	options := &webp.Options{Quality: float32(s.quality)}
	if err := webp.Encode(f, src, options); err != nil {
		return fmt.Errorf("webp encode failed: %w", err)
	}
	return nil
}

// UpdateImage updates image metadata and properties
func UpdateImage(imageID uint, req request.UpdateImageReq) (*model.Image, error) {
	// Get current image data
	image, err := op.GetImageByID(imageID)
	if err != nil {
		return nil, fmt.Errorf("image not found: %w", err)
	}

	// Track changes
	updated := false

	// Update fields if provided
	if req.Description != "" {
		image.Description = req.Description
		updated = true
	}

	if req.IsPublic != nil {
		image.IsPublic = *req.IsPublic
		updated = true
	}

	// Save changes if needed
	if updated {
		if err := op.UpdateImage(image); err != nil {
			return nil, fmt.Errorf("failed to update image metadata: %w", err)
		}
	}

	return image, nil
}

// DeleteImage removes an image and its files
func DeleteImage(imageID uint) (*response.ImageDeleteResponse, error) {
	// Get image data
	image, err := op.GetImageByID(imageID)
	if err != nil {
		return nil, fmt.Errorf("image not found: %w", err)
	}

	// Delete from database
	if err := op.DeleteImage(imageID); err != nil {
		return nil, fmt.Errorf("failed to delete image from database: %w", err)
	}

	// Delete files asynchronously
	go func() {
		if err := utils.RemoveFile(image.Path); err != nil {
			log.Printf("Warning: failed to delete image file: %v", err)
		}

		if image.ThumbnailPath != "" {
			if err := utils.RemoveFile(image.ThumbnailPath); err != nil {
				log.Printf("Warning: failed to delete thumbnail: %v", err)
			}
		}
	}()

	// Return success response
	return &response.ImageDeleteResponse{
		Message: "Image deleted successfully",
		ID:      imageID,
	}, nil
}

// RemoveTagFromImage removes a tag from an image
func RemoveTagFromImage(imageID, tagID uint) (*model.Tag, error) {
	// Remove the tag
	tag, err := op.RemoveTagFromImage(imageID, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove tag: %w", err)
	}

	return tag, err
}

// isAllowedExtension checks if the file extension is allowed
func (s *ImageService) isAllowedExtension(ext string) bool {
	ext = strings.ToLower(ext)
	for _, allowed := range s.allowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}

// getImageInfo extracts image dimensions and format
func (s *ImageService) getImageInfo(r io.Reader) (image.Image, string, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

// GetWebPPath 根据原图路径生成WebP格式图片的路径
func GetWebPPath(imagePath string) string {
	ext := filepath.Ext(imagePath)
	return strings.TrimSuffix(imagePath, ext) + ".webp"
}

// GetThumbnailPath returns the thumbnail path for an original path
func GetThumbnailPath(originalPath string) string {
	dir := filepath.Dir(originalPath)
	filename := filepath.Base(originalPath)
	thumbnailDir := filepath.Join(dir, "thumbnails")
	return filepath.Join(thumbnailDir, filename)
}
