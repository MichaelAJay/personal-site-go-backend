package custom_errors

type OperationError struct {
	Msg string
}

func (e OperationError) Error() string {
	return e.Msg
}

func NewOperationError(msg string) error {
	return &OperationError{Msg: msg}
}
