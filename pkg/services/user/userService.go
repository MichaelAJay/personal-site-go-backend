package user

import (
	"github.com/MichaelAJay/personal-site-go-backend/pkg/models"
	"gorm.io/gorm"
)

type UserService struct {
	dbClient *gorm.DB
}

func NewUserService(dbClient *gorm.DB) *UserService {
	return &UserService{
		dbClient: dbClient,
	}
}

func (userService *UserService) GetUser(email string) (*models.User, error) {
	var user models.User
	result := userService.dbClient.Where("email = ?", email).First(&user)

	if result.Error != nil {
		// Required error handling for this method?
		return nil, result.Error
	}
	return &user, nil
}

func (userService *UserService) VerifyUser(email string) (string, error) {
	user, err := userService.GetUser(email)
	if err != nil {
		// Required error handling for this method?
		return "", err
	}

	result := userService.dbClient.Model(&user).Update("IsVerified", true)
	if result.Error != nil {
		// Required error handling for this method?
		return "", result.Error
	}

	return "Verified", nil
}
