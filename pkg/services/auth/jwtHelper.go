package auth

import (
	"strconv"
	"time"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
	"github.com/golang-jwt/jwt/v5"
)

func (authService *AuthService) signToken(input types.JWTClaimsInput) (string, error) {
	claims := types.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(input.UserId), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
			Issuer:    "self",
		},
		IsAdmin: input.IsAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(authService.jwtSecret)
}
