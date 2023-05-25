package repo

import (
	"context"
	"strings"
	"time"

	"greenlight/internal/movies/repoerrors"
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

	var permissions models.Permissions

	err := r.DB.SelectContext(ctx,
		&permissions,
		query,
		userID,
	)
	if err != nil {
		return permissions, err
	}

	return permissions, nil
}

/*VER ERRORRES*/
func (r permissionRepo) AddForUser(ctx context.Context, userID int64, codes ...string) error {
	query := `
        INSERT INTO users_permissions
        SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	if err != nil {
		switch {
		case strings.Contains(err.Error(), `foreign key constraint "users_permissions_user_id_fkey"`):
			return repoerrors.ErrUserPermissionsForeignKey
		default:
			return err
		}
	}

	return err
}
