package domain

type AuthRepository interface {
	ValidateCredentials(Auth) (*User, error)
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
