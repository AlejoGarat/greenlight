package repo

import (
	"context"
	"time"

	models "greenlight/internal/users/models"

	"github.com/jmoiron/sqlx"
)

type tokenRepo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *tokenRepo {
	return &tokenRepo{
		DB: db,
	}
}

func (r tokenRepo) Insert(ctx context.Context, userID int64, ttl time.Duration, scope string) (models.Token, error) {
	token, err := models.GenerateToken(userID, ttl, scope)
	if err != nil {
		return models.Token{}, err
	}

	query := `
        INSERT INTO tokens (hash, user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return models.Token{}, err
	}
	return token, nil
}

func (r tokenRepo) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	query := `
        DELETE FROM tokens 
        WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, scope, userID)
	return err
}
