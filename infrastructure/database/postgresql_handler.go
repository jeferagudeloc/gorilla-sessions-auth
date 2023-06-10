package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jeferagudeloc/gorilla-sessions-auth/application/adapter/entity"
	"github.com/jeferagudeloc/gorilla-sessions-auth/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IPostgresqlHandler interface {
	GetConnection() *gorm.DB
}

type PostgresqlHandler struct {
	db *gorm.DB
}

func NewPostgresqlHandler(c *config) (*PostgresqlHandler, error) {
	var dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.host,
		c.user,
		c.password,
		c.database,
		c.port,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		fmt.Errorf("there was an error creating database", err)
		return nil, err
	}

	err = Migration(db)
	if err != nil {
		fmt.Errorf("there was a migration error", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)

	return &PostgresqlHandler{db}, nil
}

func (postgresqlHandler PostgresqlHandler) GetUsers() ([]domain.User, error) {
	var users []entity.User
	err := postgresqlHandler.db.Model(&entity.User{}).Preload("Role").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return entity.ToUsersDomainList(users), nil

}

func (postgresqlHandler PostgresqlHandler) ValidateCredentials(auth domain.Auth) (*domain.User, error) {
	var user entity.User
	if err := postgresqlHandler.db.Preload("Company").Preload("Profile.Permissions").Preload(clause.Associations).Where("email = ?", auth.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(domain.ErrUserNotFound)
		}
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
	if err != nil {
		return nil, errors.New(domain.ErrPasswordDoesNotMatch)
	}
	return user.ToDomain(), nil
}

func (postgresqlHandler PostgresqlHandler) CreateUser(createUserRequest *domain.CreateUserRequest) (*domain.User, error) {

	existingUser := &entity.User{}
	err := postgresqlHandler.db.Where("email = ?", createUserRequest.Email).First(existingUser).Error
	if err == nil {
		return nil, fmt.Errorf(domain.ErrConflictUser)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:        uuid.New(),
		Name:      createUserRequest.Name,
		LastName:  createUserRequest.LastName,
		Email:     createUserRequest.Email,
		Status:    createUserRequest.Status,
		Password:  string(hashedPassword),
		ProfileID: createUserRequest.ProfileID,
		CompanyID: createUserRequest.CompanyID,
	}

	if err := postgresqlHandler.db.
		Create(user).Error; err != nil {
		return nil, err
	}

	createdUser := &entity.User{}
	if err := postgresqlHandler.db.Preload("Company").Preload("Profile.Permissions").First(createdUser, user.ID).Error; err != nil {
		return nil, err
	}

	return createdUser.ToDomain(), nil
}
