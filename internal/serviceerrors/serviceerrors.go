package serviceerrors

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrNoRows         = errors.New("sql: no rows in result set")
	ErrDuplicateEmail = errors.New("duplicate email")
)
