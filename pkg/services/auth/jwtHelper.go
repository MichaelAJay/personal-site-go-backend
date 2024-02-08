package auth

import (
	"fmt"
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

func (authService *AuthService) ParseWithClaims(token string) (*types.JwtClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &types.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return authService.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*types.JwtClaims); ok && parsedToken.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token or cannot convert claims")
	}
}
