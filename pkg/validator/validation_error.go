package validator

import (
	"fmt"
	"strings"
)

type FieldError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

func (fe *FieldError) Error() string {
	return fmt.Sprintf("Field: %s, Value: %v, Message: %s", fe.Field, fe.Value, fe.Message)
}

type FieldErrors []FieldError

func (fe FieldErrors) Error() string {
	if len(fe) == 0 {
		return ""
	}

	var errMsgs []string
	for _, err := range fe {
		errMsgs = append(errMsgs, err.Error())
	}
	return strings.Join(errMsgs, "; ")
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation error"
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Errors: []FieldError{
			{Field: field, Message: message},
		},
	}
}

func (e *ValidationError) Append(field, message string) {
	e.Errors = append(e.Errors, FieldError{
		Field:   field,
		Message: message,
	})
}

func (e *ValidationError) HasError() bool {
	return len(e.Errors) > 0
}

func IsValidationError(err error) (*ValidationError, bool) {
	ve, ok := err.(*ValidationError)
	return ve, ok
}
