package service

import (
	"context"
	"errors"
	"time"

	"greenlight/internal/users/models"
	"greenlight/internal/users/repoerrors"
	"greenlight/internal/users/serviceerrors"
	"greenlight/pkg/jsonlog"
	"greenlight/pkg/mailer"
	"greenlight/pkg/taskutils"
)

type userService struct {
	repo               UserRepo
	permissionsService PermissionsService
	tokensRepo         TokensRepo
	logger             *jsonlog.Logger
	mailer             mailer.Mailer
}
type TokensRepo interface {
	Insert(ctx context.Context, userID int64, ttl time.Duration, scope string) (models.Token, error)
	DeleteAllForUser(ctx context.Context, scope string, userID int64) error
}

type UserRepo interface {
	Insert(ctx context.Context, user models.User) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	GetForToken(ctx context.Context, tokenScope string, tokenPlaintext string) (models.User, error)
}

type PermissionsService interface {
	AddForUser(ctx context.Context, userID int64, codes ...string) error
}

func NewUserService(repo UserRepo, tokensRepo TokensRepo, logger *jsonlog.Logger, mailer mailer.Mailer, permissionsService PermissionsService) *userService {
	return &userService{
		repo:               repo,
		tokensRepo:         tokensRepo,
		mailer:             mailer,
		logger:             logger,
		permissionsService: permissionsService,
	}
}

func (s userService) AddUser(ctx context.Context, user models.User) (models.User, error) {
	user, err := s.repo.Insert(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrDuplicateEmail):
			return models.User{}, serviceerrors.ErrDuplicatedEmail
		}
		return models.User{}, err
	}

	token, err := s.tokensRepo.Insert(ctx, user.ID, 3*24*time.Hour, models.ScopeActivation)
	if err != nil {
		return models.User{}, err
	}

	err = s.permissionsService.AddForUser(ctx, user.ID, "movies:read")
	if err != nil {
		return models.User{}, err
	}

	go taskutils.BackgroundTask(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = s.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			s.logger.PrintError(err, nil)
			return
		}
	})

	return user, nil
}

func (s userService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNoRows):
			return models.User{}, serviceerrors.ErrNoUserFound

		default:
			return models.User{}, err

		}
	}

	return user, nil
}

func (s userService) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	user, err := s.repo.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrRecordNotFound):
			return models.User{}, serviceerrors.ErrNoUserFound
		default:
			return models.User{}, err
		}
	}

	return user, nil
}

func (s userService) GetForToken(ctx context.Context, tokenScope string, tokenPlaintext string) (models.User, error) {
	user, err := s.repo.GetForToken(ctx, tokenScope, tokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrNoRows):
			return models.User{}, serviceerrors.ErrNoTokenFound
		default:
			return models.User{}, err
		}
	}

	return user, nil
}
