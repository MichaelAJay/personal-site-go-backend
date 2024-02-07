package routes

import (
	"net/http"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/auth"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
	"github.com/gin-gonic/gin"
)

func SignUp(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form types.SignUpRequestBody

		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, isAdmin, err := authService.SignUp(form)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		setAuthCookie(c, token)

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "isAdmin": isAdmin})
	}
}

func SignIn(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form types.SignInRequestBody

		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, isAdmin, err := authService.SignIn(form)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		setAuthCookie(c, token)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "isAdmin": isAdmin})
	}
}

func setAuthCookie(c *gin.Context, token string) {
	c.SetCookie(
		"at",
		token,
		3600*3,
		"/",
		"",
		false, // Secure: whether to use HTTPS only
		true,  // HttpOnly: JS should not be able to access
	)
}
