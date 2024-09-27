package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/language"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate = newValidator()

type ParamError struct {
	Param   string
	Message string
}

func (pe ParamError) Error() string {
	return fmt.Sprintf("%s: %s", strings.ToLower(pe.Param), pe.Message)
}

func newValidator() *validator.Validate {
	vv := validator.New()

	vv.RegisterValidation("mustbenumericalboolean", func(fl validator.FieldLevel) bool {
		nn := fl.Field().Int()

		return nn == 1 || nn == 0 || nn == -1
	})

	vv.RegisterValidation("mustbesupportedlanguage", func(fl validator.FieldLevel) bool {
		lang := fl.Field().String()

		if lang != string(language.EN) && lang != string(language.PT) && lang != string(language.RU) {
			return false
		}

		return true
	})

	vv.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`).Match([]byte(fl.Field().String()))
	})

	vv.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 7 {
			return false
		}

		if len(password) > 48 {
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

func ValidateErrorMessage(ctx *gin.Context, fe validator.FieldError) ParamError {
	dict := translations.GetTranslation(ctx)
	field := strings.ToLower(fe.Field())

	if fe.Tag() == "min" {
		return ParamError{
			Param:   field,
			Message: fmt.Sprintf(dict.ValidatorMinChars, fe.Param()),
		}
	}

	if fe.Tag() == "max" {
		return ParamError{
			Param:   field,
			Message: fmt.Sprintf(dict.ValidatorMaxChars, fe.Param()),
		}
	}

	switch fe.Tag() {
	case "required":
		return ParamError{
			Param:   field,
			Message: dict.ValidatorRequired,
		}
	case "mustbesupportedlanguage":
		return ParamError{
			Param:   field,
			Message: dict.ValidatorMustBeSupportedLanguage,
		}
	case "mustbenumericalboolean":
		return ParamError{
			Param:   field,
			Message: dict.ValidatorMustBeNumericalBoolean,
		}
	case "username":
		return ParamError{
			Param:   field,
			Message: dict.ValidatorIncorrectUsernameFormat,
		}
	case "email":
		return ParamError{
			Param:   field,
			Message: dict.ValidatorIncorrectEmailFormat,
		}
	case "password":
		return ParamError{
			Param:   field,
			Message: dict.ValidatorIncorrectPasswordFormat,
		}
	default:
		return ParamError{
			Param:   field,
			Message: fe.Error(),
		}
	}
}
