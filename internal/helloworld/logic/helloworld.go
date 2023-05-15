package logic

import (
	"context"
	"errors"
	"fmt"

	"greenlight/internal/helloworld/logicerrors"
	"greenlight/internal/helloworld/repositoryerrors"
	"greenlight/internal/models"
)

type helloWorldLogic struct {
	repo HelloWorldRepository
}

type HelloWorldRepository interface {
	SaveGreetedUser(context.Context, *models.User) error
	GetUser(ctx context.Context, name string) (models.User, error)
	GetAllGreetedUsers(context.Context) ([]models.User, error)
}

func NewHelloWorldLogic(repo HelloWorldRepository) *helloWorldLogic {
	return &helloWorldLogic{
		repo: repo,
	}
}

func (l helloWorldLogic) Greet(ctx context.Context, user *models.User, saveUser bool) (string, error) {
	if saveUser {
		err := l.repo.SaveGreetedUser(ctx, user)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("Hello, %s!", user.Name), nil
}

func (l helloWorldLogic) GetUserByName(ctx context.Context, name string) (models.User, error) {
	user, err := l.repo.GetUser(ctx, name)
	if err != nil {
		switch {
		case errors.Is(err, repositoryerrors.ErrRecordNotFound):
			return models.User{}, logicerrors.ErrUserDoesNotExist
		default:
			return models.User{}, err
		}
	}

	return user, nil
}

func (l helloWorldLogic) ListUsers(ctx context.Context) ([]models.User, error) {
	return l.repo.GetAllGreetedUsers(ctx)
}
