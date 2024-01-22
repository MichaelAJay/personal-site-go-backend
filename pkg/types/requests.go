package types

type ContactFormRequestBody struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required"`
}

type ToggleMessageReadStatusRequestBody struct {
	IsRead *bool `json:"isRead" binding:"required"`
}
