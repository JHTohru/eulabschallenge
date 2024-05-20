package product

import "strings"

// Input represents the product data a user can define.
type Input struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Validate validates the input and its data. A Validate call into an invalid
// input will return ErrNilInput, ErrInvalidName, or ErrInvalidDescription.
func (in *Input) Validate() error {
	if in == nil {
		return ErrNilInput
	}
	if strings.TrimSpace(in.Name) == "" {
		return ErrInvalidName
	}
	if strings.TrimSpace(in.Description) == "" {
		return ErrInvalidDescription
	}
	return nil
}
