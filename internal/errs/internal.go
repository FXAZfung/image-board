package errs

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid")
	ErrExist    = errors.New("exist")
)
