package product

import "fmt"

type ValidationError struct {
	msg string
}

func (ve *ValidationError) Error() string {
	return ve.msg
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{msg}
}

var (
	ErrNilInput           = NewValidationError("product input is nil")
	ErrInvalidName        = NewValidationError("product name is invalid")
	ErrInvalidDescription = NewValidationError("product description is invalid")
	ErrInvalidPageSize    = NewValidationError("product list page size is invalid")
	ErrInvalidPageNumber  = NewValidationError("product list page number is invalid")
	ErrNotFound           = fmt.Errorf("product not found")
)
