package serviceerrors

import "errors"

var (
	ErrNoMovieFound       = errors.New("no movie found")
	ErrNoMoviesFound      = errors.New("no movies found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrMovieTitleRequired = errors.New("movie title must be provided")
	ErrMovieYearRequired  = errors.New("movie year must be provided")
)
