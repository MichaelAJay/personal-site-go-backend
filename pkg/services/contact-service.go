package services

import (
	"fmt"
	"log"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/errors"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/models"
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

	contact := models.Contact{
		Name:    form.Name,
		Email:   form.Email,
		Message: form.Message,
	}

	dbClient := db

	result := dbClient.Create(&contact)
	if result.Error != nil {
		log.Fatalf("Error creating contact record: %v", result.Error)
	}

	log.Printf("Created contact record with ID: %d", contact.ID)
}

func (s *ContactService) GetUnreadForms() ([]types.UnreadContactForm, error) {
	var contacts []types.UnreadContactForm

	dbClient := db
	result := dbClient.Model(&models.Contact{}).Select("ID", "Name", "Email", "CreatedAt").Where("Is_Read = ?", false).Find(&contacts)
	if result.Error != nil {
		return nil, result.Error
	}
	return contacts, nil
}

func (s *ContactService) GetMessage(id uint) (models.Contact, error) {
	var message models.Contact

	dbClient := db
	result := dbClient.Where("ID = ?", id).First(&message)
	if result.Error != nil {
		return models.Contact{}, result.Error
	}

	// Retrieve the message
	if err := dbClient.Where("ID =?", id).First(&message).Error; err != nil {
		return models.Contact{}, err
	}

	// if err := dbClient.Model(&message).Update("IsRead", true).Error; err != nil {
	// 	log.Printf("Failed to update IsRead: %v", err)
	// }

	_, err := s.ToggleMessageReadStatus(message.ID, true)
	if err != nil {
		log.Printf("Failed to update IsRead: %v", err)
	}

	return message, nil
}

func (s *ContactService) ToggleMessageReadStatus(id uint, isRead bool) (string, error) {
	dbClient := db
	result := dbClient.Model(&models.Contact{}).Where("ID = ?", id).Update("IsRead", isRead)
	if result.Error != nil {
		return "Failure", result.Error
	}

	if result.RowsAffected == 0 {
		return "No Record Updated", errors.NotFoundError{Msg: "No record found with the given ID"}
	}

	return "Success", nil
}
