package service

import (
	"context"
	"errors"

	"greenlight/internal/repoerrors"
	"greenlight/internal/serviceerrors"
	"greenlight/internal/users/models"
)

type userService struct {
	repo UserRepo
}
type UserRepo interface {
	Insert(ctx context.Context, user models.User) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
}

func NewUserService(repo UserRepo) *userService {
	return &userService{
		repo: repo,
	}
}

func (s userService) AddUser(ctx context.Context, user models.User) (models.User, error) {
	user, err := s.repo.Insert(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrDuplicateEmail):
			return models.User{}, serviceerrors.ErrDuplicateEmail
		}
		return models.User{}, err
	}

	return user, nil
}

func (s userService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s userService) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	user, err := s.repo.Update(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
