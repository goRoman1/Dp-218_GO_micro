package validation

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// SignUpUserRequest takes user signup request data
type SignUpUserRequest struct {
	LoginEmail  string
	UserName    string
	UserSurname string
	Password    string
}

// Validate validates signupuser request data
func (su SignUpUserRequest) Validate() error {

	return validation.ValidateStruct(&su,
		validation.Field(&su.LoginEmail, validation.Required, is.EmailFormat, validation.Length(3, 100)),
		validation.Field(&su.UserName, validation.Required, validation.Length(3, 100)),
		validation.Field(&su.UserSurname, validation.Required, validation.Length(3, 100)),
		validation.Field(&su.Password, validation.Required, validation.Length(3, 10)),
	)
}

// SignInUserRequest takes user signin requesrt data
type SignInUserRequest struct {
	LoginEmail string
	Password   string
}

// Validate validates signinuser request data
func (su SignInUserRequest) Validate() error {

	return validation.ValidateStruct(&su,
		validation.Field(&su.LoginEmail, validation.Required, is.EmailFormat, validation.Length(3, 100)),
		validation.Field(&su.Password, validation.Required, validation.Length(3, 10)),
	)
}
