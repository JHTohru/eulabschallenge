package product

import "strings"

type Input struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

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
