package serviceerrors

import "errors"

var (
	ErrEditConflict              = errors.New("edit conflict")
	ErrDuplicateEmail            = errors.New("duplicate email")
	ErrEmailRequired             = errors.New("email required")
	ErrPswRequired               = errors.New("password required")
	ErrUserNotFound              = errors.New("user not found")
	ErrTokenNotFound             = errors.New("token not found")
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
)
