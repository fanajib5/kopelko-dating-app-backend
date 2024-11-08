package utils

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func GetErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Please provide a valid email address"
	case "min":
		return "The field must be at least " + fe.Param() + " characters long"
	case "gte":
		return "The field must be greater than or equal to " + fe.Param()
	case "oneof":
		return "The field must be one of the following values: " + fe.Param()
	default:
		return "Invalid value"
	}
}
