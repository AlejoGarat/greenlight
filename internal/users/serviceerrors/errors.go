package serviceerrors

import "errors"

var (
	ErrNoUserFound               = errors.New("no user found")
	ErrNoTokenFound              = errors.New("no token found")
	ErrDuplicatedEmail           = errors.New("duplicated email")
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
	ErrEditConflict              = errors.New("edit conflict")
)
