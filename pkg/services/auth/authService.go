package auth

import (
	"errors"
	"log"
	"strings"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/custom_errors"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/models"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/user"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	dbClient    *gorm.DB
	userService *user.UserService
	jwtSecret   []byte
}

func NewAuthService(dbClient *gorm.DB, userService *user.UserService, secret string) (*AuthService, error) {
	if secret == "" {
		return nil, errors.New("JWT_SECRET not found")
	}

	return &AuthService{
		dbClient:    dbClient,
		userService: userService,
		jwtSecret:   []byte(secret),
	}, nil
}

func (authService *AuthService) SignUp(form types.SignUpRequestBody) (string, error) {
	// Check for existing user
	var existingUser models.User
	result := authService.dbClient.Where("Email = ?", form.Email).First(&existingUser)
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

	createUserResult := authService.dbClient.Create(&user)
	if createUserResult.Error != nil {
		log.Printf("Error creating user record: %v", createUserResult.Error)
		return "", createUserResult.Error
	}

	token, err := authService.generateJWT(user.ID, false)
	return token, err
}

func (authService *AuthService) SignIn(input types.SignInRequestBody) (string, error) {
	user, err := authService.userService.GetUser(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Obfuscate record not found
			return "", custom_errors.NewOperationError("Operation error")
		}
		return "", err
	}

	// User is verified
	if !user.IsVerified {
		return "", custom_errors.NewAccountNotVerifiedError("Account not verified")
	}

	// Verify password
	isValid := verifyPassword(input.Password, user.Hashedpassword)
	if !isValid {
		return "", custom_errors.NewOperationError("Login not successful")
	}

	// Return JWT
	token, err := authService.generateJWT(user.ID, user.IsAdmin)
	return token, err
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
