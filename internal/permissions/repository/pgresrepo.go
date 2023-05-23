package repo

import (
	"context"
	"time"

	"greenlight/internal/permissions/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type permissionRepo struct {
	DB *sqlx.DB
}

func NewPermissionsRepo(db *sqlx.DB) *permissionRepo {
	return &permissionRepo{
		DB: db,
	}
}

func (r permissionRepo) GetAllForUser(ctx context.Context, userID int64) (models.Permissions, error) {
	query := `
        SELECT permissions.code
        FROM permissions
        INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
        INNER JOIN users ON users_permissions.user_id = users.id
        WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions models.Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r permissionRepo) AddForUser(ctx context.Context, userID int64, codes ...string) error {
	query := `
        INSERT INTO users_permissions
        SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	return err
}
