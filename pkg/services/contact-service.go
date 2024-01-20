package services

import (
	"fmt"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
)

type ContactService struct {
	// Add dependencies
}

func NewContactService() *ContactService {
	return &ContactService{}
}

func (s *ContactService) ProcessForm(form types.ContactFormRequestBody) {
	// This is working as intended
	fmt.Printf("Name: %s, Email: %s, Message: %s\n", form.Name, form.Email, form.Message)
}
