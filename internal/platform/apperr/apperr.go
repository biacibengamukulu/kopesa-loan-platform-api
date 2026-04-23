package apperr

import "github.com/biangacila/kopesa-loan-platform-api/internal/shared"

type Error struct {
	Status  int
	Code    string
	Message string
	Details []shared.FieldError
}

func (e *Error) Error() string {
	return e.Message
}

func New(status int, code, message string, details ...shared.FieldError) *Error {
	return &Error{Status: status, Code: code, Message: message, Details: details}
}
