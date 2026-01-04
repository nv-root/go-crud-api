package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationErrType struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationErrType {
	errs := []ValidationErrType{}
	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		var message string

		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "min":
			message = fmt.Sprintf("%s must have at least %s characters", field, e.Param())
		case "oneof":
			message = fmt.Sprintf("%s must be one of: %s", field, e.Param())
		case "gte":
			message = fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
		case "lte":
			message = fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
		default:
			message = fmt.Sprintf("%s has an invalid value", field)
		}

		errs = append(errs, ValidationErrType{Path: field, Message: message})
	}
	return errs
}
