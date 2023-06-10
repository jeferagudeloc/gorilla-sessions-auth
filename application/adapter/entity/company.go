package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid"`
	Name string
}

func (u *Company) TableName() string {
	return "company"
}
