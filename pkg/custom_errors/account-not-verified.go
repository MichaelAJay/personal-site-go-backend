package custom_errors

type AccountNotVerifiedError struct {
	Msg string
}

func (e AccountNotVerifiedError) Error() string {
	return e.Msg
}

func NewAccountNotVerifiedError(msg string) error {
	return &AccountNotVerifiedError{Msg: msg}
}
