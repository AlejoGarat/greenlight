package repoerrors

import (
	"errors"
)

var (
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEmailRequired  = errors.New("email required")
	ErrPswRequired    = errors.New("password required")
	ErrUserNotFound   = errors.New("user not found")
	ErrTokenNotFound  = errors.New("token not found")
	ErrUserIdRequired = errors.New("user id required")
)
