package repo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"strings"
	"time"

	"greenlight/internal/users/models"
	"greenlight/internal/users/repoerrors"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
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
		case strings.Contains(err.Error(), `unique constraint "users_email_key"`):
			return models.User{}, repoerrors.ErrDuplicateEmail
		case strings.Contains(err.Error(), `null value in column "email"`):
			return models.User{}, repoerrors.ErrEmailRequired
		case strings.Contains(err.Error(), `null value in column "password_hash"`):
			return models.User{}, repoerrors.ErrPswRequired
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

	err := r.DB.GetContext(ctx, &user, query, user.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repoerrors.ErrUserNotFound
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
		user.Password,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.GetContext(ctx, &user, query, args...)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), `unique constraint "users_email_key"`):
			return models.User{}, repoerrors.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repoerrors.ErrEditConflict
		case strings.Contains(err.Error(), `null value in column "email"`):
			return models.User{}, repoerrors.ErrEmailRequired
		case strings.Contains(err.Error(), `null value in column "password_hash"`):
			return models.User{}, repoerrors.ErrPswRequired
		default:
			return models.User{}, err
		}
	}

	return user, nil
}

func (r userRepo) GetForToken(ctx context.Context, tokenScope string, tokenPlaintext string) (models.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT u.id, u.created_at, u.name, u.email, u.password_hash, u.activated, u.version
        FROM users AS u
        INNER JOIN tokens as t
        ON u.id = t.user_id
        WHERE t.hash = $1
        AND t.scope = $2 
        AND t.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user models.User

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.GetContext(ctx, &user, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repoerrors.ErrTokenNotFound
		default:
			return models.User{}, err
		}
	}

	return user, nil
}
