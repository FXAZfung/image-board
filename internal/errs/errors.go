package errs

import (
	"errors"
	"fmt"
	pkgerr "github.com/pkg/errors"
)

var (
	NotImplement = errors.New("not implement")
	NotSupport   = errors.New("not support")
	RelativePath = errors.New("access using relative path is not allowed")

	MoveBetweenTwoStorages = errors.New("can't move files between two storages, try to copy")
	UploadNotSupported     = errors.New("upload not supported")

	MetaNotFound     = errors.New("meta not found")
	StorageNotFound  = errors.New("storage not found")
	StreamIncomplete = errors.New("upload/download stream incomplete, possible network issue")
	StreamPeekFail   = errors.New("StreamPeekFail")
)

// Basic image operation errors
var (
	ImageNotFound  = errors.New("image not found")
	ErrImageList   = errors.New("failed to list images")
	ImageSaveError = errors.New("failed to save image")
	ErrImageDelete = errors.New("failed to delete image")
	ErrImageCount  = errors.New("failed to count images")
	ErrImageUpdate = errors.New("failed to update image")
)

// File validation errors
var (
	ErrNoFileProvided    = errors.New("no file provided")
	ErrInvalidFileType   = errors.New("invalid file type, only images are accepted")
	ErrInvalidFileExt    = errors.New("invalid file extension")
	ErrFileTooLarge      = errors.New("file size exceeds maximum limit")
	ErrCorruptedFile     = errors.New("corrupted or invalid image file")
	ErrFileNameCollision = errors.New("file name collision detected")
)

// Storage errors
var (
	ErrStorageCreate     = errors.New("failed to create storage directory")
	ErrFileRead          = errors.New("failed to read file data")
	ErrFileWrite         = errors.New("failed to write file data")
	ErrStorageQuota      = errors.New("storage quota exceeded")
	ErrThumbnailGenerate = errors.New("failed to generate thumbnail")
)

// Permission errors
var (
	ErrImageAccess       = errors.New("insufficient permissions to access image")
	ErrImageNotOwned     = errors.New("user does not own this image")
	ErrImageModifyDenied = errors.New("not authorized to modify this image")
)

// Tag related errors
var (
	ErrTagsEmpty   = errors.New("no tags provided")
	ErrTagNotFound = errors.New("tag not found")
	ErrMainTagSet  = errors.New("failed to set main tag")
	ErrTagAdd      = errors.New("failed to add tags to image")
	ErrTagRemove   = errors.New("failed to remove tag from image")
	ErrTooManyTags = errors.New("maximum number of tags exceeded")
)

// Rate limiting errors
var (
	ErrTooManyRequests = errors.New("too many requests, please try again later")
	ErrImageBatchLimit = errors.New("batch operation limit exceeded")
)

// Duplication errors
var (
	ErrDuplicateImage = errors.New("image already exists in the system")
	ErrDuplicateHash  = errors.New("image with identical content already exists")
)

// User errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserCreationFailed = errors.New("failed to create user")
	ErrUserUpdateFailed   = errors.New("failed to update user")
	ErrUserDeletionFailed = errors.New("failed to delete user")
	ErrUserNotAuthorized  = errors.New("user not authorized")
)

// String errors
var (
	ErrEmptyString = errors.New("string is empty")
)

// NewErr wrap constant error with an extra message
// use errors.Is(err1, StorageNotFound) to check if err belongs to any internal error
func NewErr(err error, format string, a ...any) error {
	return fmt.Errorf("%w; %s", err, fmt.Sprintf(format, a...))
}

func IsNotFoundError(err error) bool {
	return errors.Is(pkgerr.Cause(err), ObjectNotFound) || errors.Is(pkgerr.Cause(err), StorageNotFound)
}

func IsNotSupportError(err error) bool {
	return errors.Is(pkgerr.Cause(err), NotSupport)
}
func IsNotImplement(err error) bool {
	return errors.Is(pkgerr.Cause(err), NotImplement)
}
