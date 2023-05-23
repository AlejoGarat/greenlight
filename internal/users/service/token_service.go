package service

import (
	"context"
	"errors"
	"time"

	"greenlight/internal/movies/repoerrors"
	"greenlight/internal/users/models"
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

func (s *tokenService) Insert(ctx context.Context, userID int64, ttl time.Duration, scope string) (models.Token, error) {
	token, err := s.repo.Insert(ctx, userID, ttl, scope)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNoRows):
			return models.Token{}, serviceerrors.ErrNoUserFound
		default:
			return models.Token{}, err
		}
	}
	return token, nil
}
