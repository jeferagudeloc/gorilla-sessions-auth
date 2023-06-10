package domain

import "github.com/google/uuid"

const (
	ErrUserNotFound         = "user not found"
	ErrConflictUser         = "err conflict user"
	ErrPasswordDoesNotMatch = "password does not match"
)

type (
	UserRepository interface {
		GetUsers() ([]User, error)
		CreateUser(createUserRequest *CreateUserRequest) (*User, error)
	}

	CreateUserRequest struct {
		Name      string    `json:"name"`
		LastName  string    `json:"lastname"`
		Email     string    `json:"email"`
		Status    string    `json:"status"`
		Password  string    `json:"password"`
		ProfileID uuid.UUID `json:"profileId"`
		CompanyID uuid.UUID `json:"companyId"`
	}

	User struct {
		Name     string  `json:"name"`
		LastName string  `json:"lastname"`
		Email    string  `json:"email"`
		Status   string  `json:"status"`
		Profile  Profile `json:"profile"`
		Company  Company `json:"company"`
	}

	Profile struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}

	Company struct {
		Name string `json:"name"`
	}
)
