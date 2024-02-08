package types

import "github.com/golang-jwt/jwt/v5"

type JWTClaimsInput struct {
	UserId  uint
	IsAdmin bool
}

type JwtClaims struct {
	jwt.RegisteredClaims
	IsAdmin bool
}
