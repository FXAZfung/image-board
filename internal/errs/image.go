package errs

import "errors"

var (
	ImageNotFound  = errors.New("image not found")
	ErrImageList   = errors.New("image list error")
	ImageSaveError = errors.New("image save error")
	ErrImageDelete = errors.New("image delete error")
	ErrImageCount  = errors.New("image count error")
)
