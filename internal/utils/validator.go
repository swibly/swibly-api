package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func initValidator() {
	Validate = validator.New()

	Validate.RegisterValidation("haveSpecial", func(fl validator.FieldLevel) bool {
		// `\W` = `[^a-zA-Z0-9_]`
		return regexp.MustCompile(`[\W_]`).MatchString(fl.Field().String())
	})

	Validate.RegisterValidation("haveNumeric", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`[0-9]`).MatchString(fl.Field().String())
	})
}
