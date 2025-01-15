package errs

import "errors"

var (
	ErrEmptyUsername = errors.New("empty username")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExist     = errors.New("user exist")
	ErrUserPassword  = errors.New("password error")
	ErrUserAuth      = errors.New("auth error")
	ErrUserRegister  = errors.New("register error")
	ErrUserToken     = errors.New("token error")
)
