package usecase

import (
	"context"

	"github.com/jeferagudeloc/gorilla-sessions-auth/domain"
)

type (
	CreateUserUseCase interface {
		Execute(context.Context, *domain.CreateUserRequest) (*domain.User, error)
	}

	CreateUserInteractor struct {
		repo domain.UserRepository
	}
)

func NewCreateUserInteractor(
	repo domain.UserRepository,
) CreateUserUseCase {
	return CreateUserInteractor{
		repo: repo,
	}
}

func (a CreateUserInteractor) Execute(ctx context.Context, createUserRequest *domain.CreateUserRequest) (*domain.User, error) {
	createdUser, err := a.repo.CreateUser(createUserRequest)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}
