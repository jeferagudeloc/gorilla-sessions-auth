package repository

import "github.com/jeferagudeloc/gorilla-sessions-auth/domain"

type PostgresSQL struct {
	db SQL
}

func NewPostgresSQL(db SQL) PostgresSQL {
	return PostgresSQL{
		db: db,
	}
}

func (m PostgresSQL) GetUsers() ([]domain.User, error) {
	return m.db.GetUsers()
}

func (m PostgresSQL) ValidateCredentials(auth domain.Auth) (*domain.User, error) {
	return m.db.ValidateCredentials(auth)
}

func (m PostgresSQL) CreateUser(createUserRequest *domain.CreateUserRequest) (*domain.User, error) {
	return m.db.CreateUser(createUserRequest)
}
