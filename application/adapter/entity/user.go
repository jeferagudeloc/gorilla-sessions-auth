package entity

import (
	"github.com/google/uuid"
	"github.com/jeferagudeloc/gorilla-sessions-auth/domain"
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		ID        uuid.UUID `gorm:"type:uuid"`
		Name      string
		LastName  string `gorm:"column:lastname" json:"lastname"`
		Email     string
		Password  string
		Status    string
		Profile   Profile `gorm:"foreignKey:ProfileID"`
		ProfileID uuid.UUID
		Company   Company `gorm:"foreignKey:CompanyID"`
		CompanyID uuid.UUID
	}
)

func (u *User) TableName() string {
	return "user"
}

func (u *User) ToDomain() *domain.User {
	permissions := make([]string, len(u.Profile.Permissions))
	for i, p := range u.Profile.Permissions {
		permissions[i] = p.Name
	}

	return &domain.User{
		Name:     u.Name,
		LastName: u.LastName,
		Email:    u.Email,
		Status:   u.Status,
		Profile: domain.Profile{
			Name:        u.Profile.Name,
			Permissions: mapPermissions(u.Profile.Permissions),
		},
		Company: domain.Company{
			Name: u.Company.Name,
		},
	}
}

func mapPermissions(entityPermissions []Permission) []string {
	output := make([]string, 0, len(entityPermissions))
	for _, p := range entityPermissions {
		output = append(output, p.Name)
	}

	return output
}

func ToUsersDomainList(users []User) []domain.User {
	domainUsers := make([]domain.User, len(users))
	for i, u := range users {
		domainUsers[i] = *u.ToDomain()
	}
	return domainUsers
}
