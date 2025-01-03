package errs

import "errors"

var (
	EmptyImageName = errors.New("image name is empty")
	EmptyImageType = errors.New("image type is empty")
	ImageNotFound  = errors.New("image not found")
	ImageSaveError = errors.New("image save error")
)
