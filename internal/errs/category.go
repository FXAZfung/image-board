package errs

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category exists")
	ErrCategoryCreate   = errors.New("failed to create category")
	ErrCategoryUpdate   = errors.New("failed to update category")
	ErrCategoryDelete   = errors.New("failed to delete category")
	ErrCategoryList     = errors.New("failed to list category")
	ErrCategoryGet      = errors.New("failed to get category")
	ErrCategoryCount    = errors.New("failed to count category")
)
