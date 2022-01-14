package validation

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// LocationRequest represents coordinates
type LocationRequest struct {
	Latitude  string
	Longitude string
}

// Validate validates locationrequest
func (cs LocationRequest) Validate() error {
	return validation.ValidateStruct(&cs,
		validation.Field(&cs.Latitude, validation.Required, is.Float),
		validation.Field(&cs.Longitude, validation.Required, is.Float),
	)
}
