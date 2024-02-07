package routes

import (
	"net/http"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/user"
	"github.com/gin-gonic/gin"
)

// TODO: Decide on
func VerifyUser(userService *user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// email := "TODO"
		email := "michael.a.jay@protonmail.com"

		message, err := userService.VerifyUser(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": message})
	}
}
