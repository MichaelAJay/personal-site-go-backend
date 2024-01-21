package routes

import (
	"net/http"
	"strconv"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/services"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
	"github.com/gin-gonic/gin"
)

// func PostContactFormHandler(c *gin.Context) {
// 	var form types.ContactFormRequestBody

// 	// This will bind the incoming JSON to form and handle validation
// 	if err := c.ShouldBindJSON(&form); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Internal logic here
// 	services.NewContactService().ProcessForm(form)

//		c.JSON(http.StatusOK, gin.H{"message": "Contact form submitted successfully"})
//	}
func PostContactFormHandler(service *services.ContactService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form types.ContactFormRequestBody

		// This will bind the incoming JSON to form and handle validation
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		service.ProcessForm(form)
		c.JSON(http.StatusOK, gin.H{"message": "Contact form submitted successfully"})
	}
}

// func GetUnreadContactFormList(c *gin.Context) {
// 	list, err := services.NewContactService().GetUnreadForms()

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving unread forms"})
// 		return
// 	}

//		c.JSON(http.StatusOK, list)
//	}
func GetUnreadContactFormListHandler(service *services.ContactService) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := service.GetUnreadForms()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving unread forms"})
			return
		}

		c.JSON(http.StatusOK, list)
	}
}

func GetMessageHandler(service *services.ContactService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		message, err := service.GetMessage(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving message"})
			return
		}

		c.JSON(http.StatusOK, message)
	}
}
