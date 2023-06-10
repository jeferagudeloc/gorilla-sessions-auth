package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Profile struct {
		gorm.Model
		ID          uuid.UUID `gorm:"type:uuid"`
		Name        string
		Permissions []Permission `gorm:"many2many:profile_permissions;"`
	}
)

func (u *Profile) TableName() string {
	return "profile"
}
