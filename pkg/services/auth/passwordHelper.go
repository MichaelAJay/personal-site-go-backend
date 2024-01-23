package auth

import "golang.org/x/crypto/bcrypt"

func verifyPassword(password string, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}
