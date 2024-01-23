package models

import (
	"gorm.io/gorm"
)

// This should NEVER be returned out of the server
type User struct {
	gorm.Model
	Firstname      string `gorm:"not null"`
	Lastname       string
	Email          string `gorm:"not null;uniqueIndex"`
	Hashedpassword []byte `gorm:"not null"`
}
