package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/custom_errors"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/contact"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/types"
	"github.com/gin-gonic/gin"
)

func PostContactFormHandler(service *contact.ContactService) gin.HandlerFunc {
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

func GetContactFormListHandler(service *contact.ContactService) gin.HandlerFunc {
	return func(c *gin.Context) {
		pgQueryParam := c.DefaultQuery("pg", "1")
		orderQueryParam := c.DefaultQuery("order", "created_at_desc")
		readQueryParam := c.DefaultQuery("read", "false")

		pgNum, err := strconv.Atoi(pgQueryParam)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "pg query param could not be converted to an integer"})
			return
		}

		getRead, err := strconv.ParseBool(readQueryParam)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "read query param could not be converted to a boolean"})
			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order param"})
			return
		}

		// Only ordering is on created_at
		orderParts := strings.Split(orderQueryParam, "_")
		if !(len(orderParts) == 3 && orderParts[0] == "created" && orderParts[1] == "at") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only ordering is on created_at"})
			return
		}

		list, err := service.GetMessages(pgNum, orderParts[2], getRead)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving unread forms"})
			return
		}

		c.JSON(http.StatusOK, list)
	}
}

func GetMessageHandler(service *contact.ContactService) gin.HandlerFunc {
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

func ToggleMessageReadStatus(service *contact.ContactService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		parsedId, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		id := uint(parsedId)

		var body types.ToggleMessageReadStatusRequestBody

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if body.IsRead == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "IsRead must be provided"})
			return
		}

		status, err := service.ToggleMessageReadStatus(id, *body.IsRead)
		if err != nil {
			if _, ok := err.(custom_errors.NotFoundError); ok {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": status})
	}
}
