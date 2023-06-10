package database

import (
	"os"

	"github.com/google/uuid"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.User{}, &entity.Profile{}, &entity.Permission{}, &entity.Company{})
	if err != nil {
		return err
	}

	if os.Getenv("LOAD_DEFAULT_DATA") == "true" {

		company := &entity.Company{
			ID:   uuid.New(),
			Name: "AZ Company",
		}
		db.Create(&company)

		profile := &entity.Profile{
			ID:   uuid.New(),
			Name: "admin",
			Permissions: []entity.Permission{{
				ID:   uuid.New(),
				Name: "read-dashboard",
			}},
		}
		db.Create(profile)

		user := &entity.User{
			ID:        uuid.New(),
			Name:      "Jeferson",
			LastName:  "Agudelo",
			Email:     "jefersonagudeloc@gmail.com",
			Status:    "active",
			CompanyID: company.ID,
			ProfileID: profile.ID,
		}

		password := []byte("root123")
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)

		db.Create(user)
	}

	return nil
}
