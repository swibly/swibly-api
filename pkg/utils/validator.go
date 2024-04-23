package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate = newValidator()

type ParamError struct {
	Param   string
	Message string
}

func (pe ParamError) Error() string {
	return fmt.Sprintf("%s: %s", pe.Param, pe.Message)
}

func newValidator() *validator.Validate {
	vv := validator.New()

	vv.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		if fl.Field().Len() > 32 {
			return false
		}

		if fl.Field().Len() < 3 {
			return false
		}

		return regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`).Match(fl.Field().Bytes())
	})

	vv.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 7 {
			return false
		}

		// \d stands for digits [0-9]
		// \W_ stands for [^a-zA-Z0-9_]

		if !regexp.MustCompile(`\d`).MatchString(password) {
			return false
		}

		if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
			return false
		}

		if !regexp.MustCompile(`[\W_]`).MatchString(password) {
			return false
		}

		if !regexp.MustCompile(`[a-z]`).MatchString(password) {
			return false
		}

		return true
	})

	return vv
}

func ValidateStruct(s any) validator.ValidationErrors {
	if err := Validate.Struct(s); err != nil {
		// We won't pass any invalid values so we can just skip that step
		return Validate.Struct(s).(validator.ValidationErrors)
	}

	return nil
}

func ValidateErrorMessage(fe validator.FieldError) ParamError {
	if fe.Tag() == "min" {
		return ParamError{
			Param:   fe.Tag(),
			Message: fmt.Sprintf("%s must have at least %s characters", fe.Tag(), fe.Param()),
		}
	}

	if fe.Tag() == "max" {
		return ParamError{
			Param:   fe.Tag(),
			Message: fmt.Sprintf("%s must have a maximum of %s characters", fe.Tag(), fe.Param()),
		}
	}

	switch fe.Tag() {
	case "required":
		return ParamError{
			Param:   fe.Tag(),
			Message: fmt.Sprintf("%s is required", fe.Tag()),
		}
	case "username":
		return ParamError{
			Param:   fe.Tag(),
			Message: "Usernames must consist only of lowercase alphanumeric (a-z & 0-9) characters",
		}
	case "email":
		return ParamError{
			Param:   fe.Tag(),
			Message: "Email format is incorrect",
		}
	case "password":
		return ParamError{
			Param:   fe.Tag(),
			Message: "The password should contain at least one uppercase letter, one lowercase letter, one special character, and one numeral",
		}
	default:
		return ParamError{
			Param:   fe.Tag(),
			Message: fe.Error(),
		}
	}
}
