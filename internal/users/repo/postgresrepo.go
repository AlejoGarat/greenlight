package repo

import (
	"context"
	"errors"
	"time"

	"greenlight/internal/repoerrors"
	"greenlight/internal/users/models"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	DB *sqlx.DB
}

func NewMovieRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		DB: db,
	}
}

func (r userRepo) Insert(ctx context.Context, user models.User) (models.User, error) {
	query := `
	INSERT INTO users (name, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.GetContext(ctx, &user, query, args...)
	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			return models.User{}, repoerrors.ErrDuplicateEmail
		default:
			return models.User{}, err
		}
	}

	return user, nil
}

func (r userRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query := `
	SELECT id, created_at, name, email, password_hash, activated, version
	FROM users
	WHERE email = $1`

	var user models.User

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.GetContext(ctx, &user, query)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNoRows):
			return models.User{}, repoerrors.ErrRecordNotFound
		default:
			return models.User{}, err
		}
	}

	return user, nil
}

func (r userRepo) Update(ctx context.Context, user models.User) (models.User, error) {
	query := `
	UPDATE users 
	SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.GetContext(ctx, &user, query, args...)
	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			return models.User{}, repoerrors.ErrDuplicateEmail
		case errors.Is(err, repoerrors.ErrNoRows):
			return models.User{}, repoerrors.ErrEditConflict
		default:
			return models.User{}, err
		}
	}

	return user, nil
}
