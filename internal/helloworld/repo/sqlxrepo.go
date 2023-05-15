package repo

import (
	"context"
	"database/sql"
	"errors"
	"greenlight/internal/helloworld/repositoryerrors"
	"greenlight/internal/models"

	"github.com/jmoiron/sqlx"
)

type sqlxRepo struct {
	db *sqlx.DB
}

func NewSqlxRepo(db *sqlx.DB) *sqlxRepo {
	return &sqlxRepo{
		db: db,
	}
}

func (r *sqlxRepo) SaveGreetedUser(ctx context.Context, user *models.User) error {

	query := `INSERT INTO users (name)
				VALUES($1)
				ON CONFLICT(name) DO NOTHING`

	row := r.db.QueryRowxContext(ctx, query, user.Name)
	if err := row.Err(); err != nil {
		return err
	}

	err := row.StructScan(user)
	if err != nil {
		return err
	}

	return nil
}

func (r *sqlxRepo) GetUser(ctx context.Context, name string) (models.User, error) {
	var user models.User

	query := `SELECT name, id, registered_at FROM users
				WHERE name = $1`

	// db.Get loads the first element into dest
	err := r.db.GetContext(ctx, &user, query, name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.User{}, repositoryerrors.ErrRecordNotFound
		default:
			return models.User{}, err
		}
	}

	return user, nil
}

func (r *sqlxRepo) GetAllGreetedUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `SELECT name, id, registered_at
				FROM users`

	// db.Select loads a slice of elements into dest
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.User{}, nil
		default:
			return []models.User{}, err
		}
	}

	return users, nil
}
