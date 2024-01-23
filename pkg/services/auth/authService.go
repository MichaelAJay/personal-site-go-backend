package auth

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/custom_errors"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/models"
	db_client "github.com/MichaelAJay/personal-site-go-backend/pkg/services/db-client"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService() (*AuthService, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set in environment")
	}

	return &AuthService{
		jwtSecret: []byte(secret),
	}, nil
}

func (authService *AuthService) SignUp(form types.SignUpRequestBody) (string, error) {
	dbClient := db_client.Db

	// Check for existing user
	var existingUser models.User
	result := dbClient.Where("Email = ?", form.Email).First(&existingUser)
	// result.Error SHOULD BE gorm.ErrRecordNotFound.
	// If it is, continue execution
	// If not
	// - If some error OTHER than gorm.ErrRecordNotFound, return it
	// - Otherwise no error - a record was found already with a matching email
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if result.Error != nil {
			return "", result.Error
		} else { // NO ERROR - User exists
			return "", custom_errors.NewOperationError("Unable to process registration")
		}
	}

	// Proceed with user creation

	// Parse name
	names := strings.Fields(form.Name)
	var firstname, lastname string
	if len(names) > 0 {
		firstname = names[0]
		if len(names) > 1 {
			lastname = strings.Join(names[1:], "")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", err

	}
	user := models.User{
		Firstname:      firstname,
		Lastname:       lastname,
		Email:          form.Email,
		Hashedpassword: hashedPassword,
	}

	createUserResult := dbClient.Create(&user)
	if createUserResult.Error != nil {
		log.Printf("Error creating user record: %v", createUserResult.Error)
		return "", createUserResult.Error
	}

	return authService.generateJWT(user.ID, false)
}

func (authService *AuthService) SignIn(input types.SignInRequestBody) (string, error) {
	dbClient := db_client.Db

	var user models.User
	result := dbClient.Where("Email = ?", input.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", custom_errors.NewOperationError("Operation error")
		} else {
			return "", result.Error
		}
	}

	// Verify password
	isValid := verifyPassword(input.Password, user.Hashedpassword)
	if !isValid {
		return "", custom_errors.NewOperationError("Login not successful")
	}

	// Return JWT
	return authService.generateJWT(user.ID, false)
}

func (authService *AuthService) generateJWT(userId uint, isAdmin bool) (string, error) {
	input := types.JWTClaimsInput{
		UserId:  userId,
		IsAdmin: isAdmin,
	}
	token, err := authService.signToken(input)
	if err != nil {
		log.Printf("Error signing credentials token: %v", err)
		return "", err
	}
	return token, nil
}
