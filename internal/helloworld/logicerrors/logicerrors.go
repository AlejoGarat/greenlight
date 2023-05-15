package logicerrors

import "errors"

var (
	ErrUserDoesNotExist = errors.New("user does not exists")
)
