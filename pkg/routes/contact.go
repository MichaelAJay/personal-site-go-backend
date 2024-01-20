package routes

import (
	"net/http"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/services"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
	"github.com/gin-gonic/gin"
)

// type ContactFormRequestBody struct {
// 	Name    string `json:"name" binding:"required"`
// 	Email   string `json:"email" binding:"required,email"`
// 	Message string `json:"message" binding:"required"`
// }

func ContactHandler(c *gin.Context) {
	var form types.ContactFormRequestBody

	// This will bind the incoming JSON to form and handle validation
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Internal logic here
	services.NewContactService().ProcessForm(form)

	c.JSON(http.StatusOK, gin.H{"message": "Contact form submitted successfully"})
}
