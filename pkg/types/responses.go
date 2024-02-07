package types

import (
	"time"

	"gorm.io/gorm"
)

type UnreadContactForm struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserReturn struct {
	gorm.Model
	Firstname string
	Lastname  string
	Email     string
}
