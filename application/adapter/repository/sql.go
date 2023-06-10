package repository

import "github.com/jeferagudeloc/gorilla-sessions-auth/domain"

type SQL interface {
	GetUsers() ([]domain.User, error)
	ValidateCredentials(domain.Auth) (*domain.User, error)
	CreateUser(createUserRequest *domain.CreateUserRequest) (*domain.User, error)
}
