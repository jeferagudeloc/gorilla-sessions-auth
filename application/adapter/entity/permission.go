package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Permission struct {
		gorm.Model
		ID   uuid.UUID `gorm:"type:uuid"`
		Name string
	}
)

func (u *Permission) TableName() string {
	return "permission"
}
