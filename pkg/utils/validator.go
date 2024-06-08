package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/language"
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

	vv.RegisterValidation("mustbesupportedlanguage", func(fl validator.FieldLevel) bool {
		lang := fl.Field().String()

		if lang != string(language.EN) && lang != string(language.PT) && lang != string(language.RU) {
			return false
		}

		return true
	})

	vv.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		if fl.Field().Len() > 32 {
			return false
		}

		if fl.Field().Len() < 3 {
			return false
		}

		return regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`).Match([]byte(fl.Field().String()))
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
			Param:   fe.Field(),
			Message: fmt.Sprintf("%s must have at least %s characters", fe.Field(), fe.Param()),
		}
	}

	if fe.Tag() == "max" {
		return ParamError{
			Param:   fe.Field(),
			Message: fmt.Sprintf("%s must have a maximum of %s characters", fe.Field(), fe.Param()),
		}
	}

	switch fe.Tag() {
	case "required":
		return ParamError{
			Param:   fe.Field(),
			Message: fmt.Sprintf("%s is required", fe.Field()),
		}
	case "mustbesupportedlanguage":
		return ParamError{
			Param:   fe.Field(),
			Message: fmt.Sprintf("%s must be %s", fe.Field(), strings.Join(language.ArrayString, ",")),
		}
	case "username":
		return ParamError{
			Param:   fe.Field(),
			Message: "Usernames must consist only of lowercase alphanumeric (a-z & 0-9) characters",
		}
	case "email":
		return ParamError{
			Param:   fe.Field(),
			Message: "Email format is incorrect",
		}
	case "password":
		return ParamError{
			Param:   fe.Field(),
			Message: "The password should contain at least one uppercase letter, one lowercase letter, one special character, and one numeral",
		}
	default:
		return ParamError{
			Param:   fe.Field(),
			Message: fe.Error(),
		}
	}
}
