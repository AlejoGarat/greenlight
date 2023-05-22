package serviceerrors

import "errors"

var (
	ErrNoMovieFound = errors.New("no movie found")
	ErrEditConflict = errors.New("edit conflict")
)
