package services

import (
	"fmt"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/routes"
)

type ContactForm struct {
	routes.ContactFormRequestBody
}

type ContactService struct {
	// Add dependencies
}

func NewContactService() *ContactService {
	return &ContactService{}
}

func (s *ContactService) ProcessForm(form ContactForm) {
	fmt.Printf("Name: %s, Email: %s, Message: %s\n", form.Name, form.Email, form.Message)
}
