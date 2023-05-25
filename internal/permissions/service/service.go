package service

import (
	"context"
	"errors"

	"greenlight/internal/permissions/models"
	"greenlight/internal/users/repoerrors"
	"greenlight/internal/users/serviceerrors"
	"greenlight/pkg/jsonlog"
)

type permissionsService struct {
	repo   PermissionsRepo
	logger *jsonlog.Logger
}

type PermissionsRepo interface {
	AddForUser(ctx context.Context, userID int64, codes ...string) error
	GetAllForUser(ctx context.Context, userID int64) (models.Permissions, error)
}

func NewPermissionsService(repo PermissionsRepo, logger *jsonlog.Logger) *permissionsService {
	return &permissionsService{
		repo:   repo,
		logger: logger,
	}
}

func (s permissionsService) AddForUser(ctx context.Context, userID int64, codes ...string) error {
	err := s.repo.AddForUser(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrUserNotFound):
			return serviceerrors.ErrUserNotFound
		default:
			return err
		}
	}

	return nil
}
