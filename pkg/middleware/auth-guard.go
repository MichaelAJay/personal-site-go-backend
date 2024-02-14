package middleware

import (
	"net/http"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/auth"
	"github.com/gin-gonic/gin"
)

func AuthGuard(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, err := authService.ParseWithClaims(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid token"})
			return
		}

		c.Set("userInfo", claims)
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	token, err := c.Cookie("at")
	if err != nil {
		return "", err
	}
	return token, nil
}
