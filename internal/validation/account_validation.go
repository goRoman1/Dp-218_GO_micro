package validation

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// CreateAccountRequest takes name and number of account
type CreateAccountRequest struct {
	Name   string
	Number string
}

// Validate validates signinuser request data
func (ca CreateAccountRequest) Validate() error {

	return validation.ValidateStruct(&ca,
		validation.Field(&ca.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&ca.Number, validation.Required, validation.Length(16, 16), is.Digit),
	)
}
