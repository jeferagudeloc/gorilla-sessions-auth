package usecase

import (
	"context"

	"github.com/jeferagudeloc/gorilla-sessions-auth/domain"
)

type (
	AuthUseCase interface {
		Execute(context.Context, domain.Auth) (*domain.User, error)
	}

	AuthInteractor struct {
		repo domain.AuthRepository
	}
)

func NewAuthInteractor(
	repo domain.AuthRepository,
) AuthUseCase {
	return AuthInteractor{
		repo: repo,
	}
}

func (a AuthInteractor) Execute(ctx context.Context, auth domain.Auth) (*domain.User, error) {
	user, err := a.repo.ValidateCredentials(auth)

	if err != nil {
		return nil, err
	}
	return user, nil
}
