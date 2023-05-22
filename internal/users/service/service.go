package service

import (
	"context"
	"errors"

	"greenlight/internal/repoerrors"
	"greenlight/internal/serviceerrors"
	"greenlight/internal/users/models"
	"greenlight/pkg/jsonlog"
	"greenlight/pkg/mailer"
)

type userService struct {
	repo   UserRepo
	logger *jsonlog.Logger
	mailer mailer.Mailer
}
type UserRepo interface {
	Insert(ctx context.Context, user models.User) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
}

func NewUserService(repo UserRepo, logger *jsonlog.Logger, mailer mailer.Mailer) *userService {
	return &userService{
		repo:   repo,
		mailer: mailer,
		logger: logger,
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

	go func() {
		err = s.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			s.logger.PrintError(err, nil)
		}
	}()

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
