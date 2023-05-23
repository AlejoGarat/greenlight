package service

import (
	"context"
	"errors"

	"greenlight/internal/movies/repoerrors"
	"greenlight/internal/users/serviceerrors"
	"greenlight/pkg/jsonlog"
)

type tokenService struct {
	repo   TokensRepo
	logger *jsonlog.Logger
}

func NewTokensService(repo TokensRepo, logger *jsonlog.Logger) *tokenService {
	return &tokenService{
		repo:   repo,
		logger: logger,
	}
}

func (s *tokenService) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {
	err := s.repo.DeleteAllForUser(ctx, scope, userID)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrRecordNotFound):
			return serviceerrors.ErrNoTokenFound
		}
		return err
	}
	return nil
}
