package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

// ImageService handles all image processing operations
type ImageService struct {
	baseDir        string
	thumbnailWidth int
	quality        int
	allowedExts    []string
}

// NewImageService creates a new image service instance
func NewImageService() *ImageService {
	return &ImageService{
		baseDir:        config.Conf.DataImage.Dir,
		thumbnailWidth: 300,
		quality:        90,
		allowedExts:    []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	}
}

// UploadImage handles the complete image upload process from handler
func UploadImage(file *multipart.FileHeader, user *model.User, req request.UploadImageReq) (*model.Image, error) {
	service := NewImageService()

	// Process the upload
	image, err := service.processUpload(file, user)
	if err != nil {
		return nil, err
	}

	// Update additional fields from the request
	if req.Description != "" {
		image.Description = req.Description
	}

	image.IsPublic = req.IsPublic

	err = op.AddTagsToImage(image.ID, req.Tags)
	if err != nil {
		return nil, err
	}

	// Save updated image metadata
	if err := op.UpdateImage(image); err != nil {
		log.Printf("Warning: failed to update image metadata: %v", err)
	}

	return image, nil
}

// processUpload handles the core upload functionality
func (s *ImageService) processUpload(file *multipart.FileHeader, user *model.User) (*model.Image, error) {
	// Validation
	if file == nil {
		return nil, errors.New("no file provided")
	}

	if !utils.IsImage(file) {
		return nil, errors.New("invalid file type, only images are accepted")
	}

	// Open and read file
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	// Generate hash for deduplication
	hash := sha256.Sum256(fileData)
	hashString := hex.EncodeToString(hash[:])

	// Check for duplicate image
	if existingImage, err := op.GetImageByHash(hashString); err == nil {
		return existingImage, nil
	}

	// Process file metadata
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	if !s.isAllowedExtension(fileExt) {
		return nil, errors.New("invalid file extension")
	}

	contentType := http.DetectContentType(fileData)
	newFileName := hashString + fileExt

	// Create storage paths using date-based directory structure
	now := time.Now()
	dirPath := filepath.Join(s.baseDir, fmt.Sprintf("%d/%02d", now.Year(), now.Month()))
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	filePath := filepath.Join(dirPath, newFileName)
	thumbnailDir := filepath.Join(dirPath, "thumbnails")
	thumbnailPath := filepath.Join(thumbnailDir, newFileName)

	// Extract image metadata
	imgInfo, _, err := s.getImageInfo(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("invalid image data: %w", err)
	}

	// Process files concurrently
	g, _ := errgroup.WithContext(context.Background())

	// Save original file
	g.Go(func() error {
		return os.WriteFile(filePath, fileData, 0644)
	})

	// Create thumbnail
	var thumbnailErr error
	g.Go(func() error {
		if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
			thumbnailErr = err
			return nil // Continue even if thumbnail dir creation fails
		}

		// Generate thumbnail
		thumbnailErr = s.createThumbnail(fileData, thumbnailPath)
		return nil // Continue even if thumbnail creation fails
	})

	// Wait for file operations to complete
	if err := g.Wait(); err != nil {
		// Clean up on error
		os.Remove(filePath)
		os.Remove(thumbnailPath)
		return nil, fmt.Errorf("file processing failed: %w", err)
	}

	// Create image model
	image := &model.Image{
		FileName:     newFileName,
		OriginalName: file.Filename,
		Hash:         hashString,
		ContentType:  contentType,
		Size:         file.Size,
		Path:         filePath,
		Width:        imgInfo.Bounds().Dx(),
		Height:       imgInfo.Bounds().Dy(),
		UserID:       user.ID,
		IsPublic:     true, // Default to public
	}

	// Set thumbnail path if successful
	if thumbnailErr == nil {
		image.ThumbnailPath = thumbnailPath
	} else {
		log.Printf("Warning: thumbnail generation failed: %v", thumbnailErr)
	}

	// Save to database
	if err := op.CreateImage(image); err != nil {
		// Clean up files on database error
		os.Remove(filePath)
		if image.ThumbnailPath != "" {
			os.Remove(image.ThumbnailPath)
		}
		return nil, fmt.Errorf("failed to save image to database: %w", err)
	}

	return image, nil
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

// AddTagsToImage adds tags to an image
func AddTagsToImage(imageID uint, tags []string) (*response.ImageTagResponse, error) {
	// Verify image exists
	if _, err := op.GetImageByID(imageID); err != nil {
		return nil, fmt.Errorf("image not found: %w", err)
	}

	// Add tags
	if err := op.AddTagsToImage(imageID, tags); err != nil {
		return nil, fmt.Errorf("failed to add tags: %w", err)
	}

	return &response.ImageTagResponse{
		ImageID: imageID,
		Success: true,
	}, nil
}

// RemoveTagFromImage removes a tag from an image
func RemoveTagFromImage(imageID, tagID uint) (*response.ImageTagResponse, error) {
	// Remove the tag
	if err := op.RemoveTagFromImage(imageID, tagID); err != nil {
		return nil, fmt.Errorf("failed to remove tag: %w", err)
	}

	return &response.ImageTagResponse{
		ImageID: imageID,
		Success: true}, nil
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

// createThumbnail generates and saves a thumbnail
func (s *ImageService) createThumbnail(imgData []byte, thumbnailPath string) error {
	// Read image
	src, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return err
	}

	// Resize keeping aspect ratio
	thumbnail := imaging.Resize(src, s.thumbnailWidth, 0, imaging.Lanczos)

	// Save thumbnail
	return imaging.Save(thumbnail, thumbnailPath, imaging.JPEGQuality(s.quality))
}

// GetThumbnailPath returns the thumbnail path for an original path
func GetThumbnailPath(originalPath string) string {
	dir := filepath.Dir(originalPath)
	filename := filepath.Base(originalPath)
	thumbnailDir := filepath.Join(dir, "thumbnails")
	return filepath.Join(thumbnailDir, filename)
}

// RegenerateMissingThumbnails creates thumbnails for images that don't have them
func RegenerateMissingThumbnails() (int, error) {
	images, _, err := op.GetImagesByPage(1, 1000) // Process in batches
	if err != nil {
		return 0, err
	}

	service := NewImageService()
	count := 0

	for _, img := range images {
		// Skip if thumbnail exists and file exists
		if img.ThumbnailPath != "" && utils.IsExist(img.ThumbnailPath) {
			continue
		}

		// Skip if original doesn't exist
		if !utils.IsExist(img.Path) {
			continue
		}

		// Create thumbnail directory
		thumbnailDir := filepath.Join(filepath.Dir(img.Path), "thumbnails")
		if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
			log.Printf("Failed to create thumbnail directory for %s: %v", img.FileName, err)
			continue
		}

		// Set thumbnail path
		thumbnailPath := filepath.Join(thumbnailDir, filepath.Base(img.Path))

		// Read original file
		fileData, err := os.ReadFile(img.Path)
		if err != nil {
			log.Printf("Failed to read original image %s: %v", img.FileName, err)
			continue
		}

		// Create thumbnail
		if err := service.createThumbnail(fileData, thumbnailPath); err != nil {
			log.Printf("Failed to create thumbnail for %s: %v", img.FileName, err)
			continue
		}

		// Update database record
		img.ThumbnailPath = thumbnailPath
		if err := op.UpdateImage(img); err != nil {
			log.Printf("Failed to update thumbnail path for %s: %v", img.FileName, err)
			continue
		}

		count++
	}

	return count, nil
}
