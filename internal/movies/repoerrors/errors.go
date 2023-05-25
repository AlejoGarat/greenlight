package repoerrors

import (
	"errors"
)

var (
	ErrEditConflict              = errors.New("edit conflict")
	ErrMovieNoFound              = errors.New("movie no found")
	ErrMoviesNoFound             = errors.New("movies no found")
	ErrMovieTitleRequired        = errors.New("movie title required")
	ErrMovieYearRequired         = errors.New("movie year required")
	ErrInvalidId                 = errors.New("invalid id")
	ErrUserPermissionsForeignKey = errors.New("user permissions foreign key")
)
