package model

import (
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type userBuildError struct {
	Param   string `json:"param"`
	Message string `json:"message"`
}

type User struct {
	gorm.Model
	Fullname string `              validate:"required,min=2,max=100"`
	Username string `gorm:"unique" validate:"required,lowercase,min=3,max=20,alphanum"`
	Email    string `gorm:"unique" validate:"required,email"`
	Password string `              validate:"required,haveSpecial,haveNumeric,min=8,max=48"`
}

func msgForTag(fe validator.FieldError) string {
	if strings.Contains(fe.Tag(), "min") {
		return fe.Field() + " is too short. It must be at least " + fe.Param() + " characters long."
	}

	if strings.Contains(fe.Tag(), "max") {
		return fe.Field() + " is too long. It must be no more than " + fe.Param() + " characters long."
	}

	switch fe.Tag() {
	case "required":
		return "Please ensure that the " + strings.ToLower(fe.Field()) + " is filled out."
	case "alphanum":
		return fe.Field() + " must contain only alphanumeric characters."
	case "lowercase":
		return fe.Field() + " must be all lowercase."
	case "email":
		return "Please enter a valid email address."
	case "haveSpecial":
		return fe.Field() + " must include at least one symbol."
	case "haveNumeric":
		return fe.Field() + " must include at least one number."
	default:
		return fe.Error()
	}
}

func (u *User) Validate() ([]userBuildError, error) {
	// [ ]: Benchmark: Test if there is any impact upon creating a new validator every time it's called
	var validate = validator.New()

	validate.RegisterValidation("haveSpecial", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`[\W_]`).MatchString(fl.Field().String())
	})

	validate.RegisterValidation("haveNumeric", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`[0-9]`).MatchString(fl.Field().String())
	})

	err := validate.Struct(u)

	var ve validator.ValidationErrors

	// Type assertion to get only validation errors
	if errors.As(err, &ve) {
		// `make` allocates memory for further implementation
		out := make([]userBuildError, len(ve))

		for i, fe := range ve {
			out[i] = userBuildError{fe.Field(), msgForTag(fe)}
		}

		return out, nil
	}

	// If `err` is `nil`, it will return `nil` anyways
	// no need to check if err != nil then return err
	return nil, err
}
