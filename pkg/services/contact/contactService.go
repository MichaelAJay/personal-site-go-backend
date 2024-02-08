package contact

import (
	"fmt"
	"log"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/custom_errors"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/models"
	db_client "github.com/MichaelAJay/personal-site-go-backend/pkg/services/db-client"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
)

type ContactService struct {
	// Add dependencies
}

func NewContactService() *ContactService {
	return &ContactService{}
}

func (s *ContactService) ProcessForm(form types.ContactFormRequestBody) (string, error) {
	// This is working as intended
	fmt.Printf("Name: %s, Email: %s, Message: %s\n", form.Name, form.Email, form.Message)

	contact := models.Contact{
		Name:    form.Name,
		Email:   form.Email,
		Message: form.Message,
	}

	dbClient := db_client.Db

	result := dbClient.Create(&contact)
	if result.Error != nil {
		log.Printf("Error creating contact record: %v", result.Error)
		return "Failure", result.Error
	}

	log.Printf("Created contact record with ID: %d", contact.ID)
	return "Success", nil
}

/*
 * createdAtOrder should be either "asc" or "desc"
 */

// TODO - set up the pagination results
func (s *ContactService) GetMessages(pgNum int, createdAtOrderDirection string, getRead bool) ([]types.UnreadContactForm, error) {
	var messages []types.UnreadContactForm

	offset := (pgNum - 1) * 10
	orderArg := fmt.Sprintf("created_at %s", createdAtOrderDirection)

	dbClient := db_client.Db
	result := dbClient.Model(&models.Contact{}).Select("ID", "Name", "Email", "CreatedAt").Where("Is_Read = ?", getRead).Find(&messages).Order(orderArg).Offset(offset).Limit(10)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (s *ContactService) GetMessage(id uint) (models.Contact, error) {
	var message models.Contact

	dbClient := db_client.Db
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
	dbClient := db_client.Db
	result := dbClient.Model(&models.Contact{}).Where("ID = ?", id).Update("IsRead", isRead)
	if result.Error != nil {
		return "Failure", result.Error
	}

	if result.RowsAffected == 0 {
		return "No Record Updated", custom_errors.NotFoundError{Msg: "No record found with the given ID"}
	}

	return "Success", nil
}
