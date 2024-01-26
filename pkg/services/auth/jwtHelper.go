package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
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

	keyPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	if keyPath == "" {
		return "", fmt.Errorf("nothing found at key path")
	}
	// keyPath looks good
	// go_dev_RSA256.key
	key, err := readPrivateKey(keyPath)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func (authService *AuthService) ParseWithClaims(token string) (*types.JwtClaims, error) {
	keyPath := os.Getenv("RSA_PUBLIC_KEY_PATH")
	if keyPath == "" {
		return nil, fmt.Errorf("nothing found at key path")
	}

	key, err := readPublicKey(keyPath)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.ParseWithClaims(token, &types.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
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

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// keyData is a []uint8

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	// block looks fine, see writeup.txt

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func readPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unknown type of public key")
	}

	return rsaPubKey, nil
}
