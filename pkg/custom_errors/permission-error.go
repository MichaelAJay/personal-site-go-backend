package custom_errors

type PermissionError struct {
	Message string
}

func (e *PermissionError) Error() string {
	return e.Message
}

func NewPermissionError(message string) *PermissionError {
	return &PermissionError{Message: message}
}
