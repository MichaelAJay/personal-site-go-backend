package types

type ContactFormRequestBody struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required"`
}

type ToggleMessageReadStatusRequestBody struct {
	IsRead *bool `json:"isRead" binding:"required"`
}

type SignUpRequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignInRequestBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreatedAtOrder string

const (
	CreatedAtDesc CreatedAtOrder = "created_at_desc"
	CreatedAtAsc  CreatedAtOrder = "created_at_asc"
)
