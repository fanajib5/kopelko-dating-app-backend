package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func getErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "gte":
		return "This field must be greater than or equal to " + fieldError.Param()
	case "lte":
		return "This field must be less than or equal to " + fieldError.Param()
	case "min":
		return "This field must be at least " + fieldError.Param() + " characters long"
	case "max":
		return "This field must be at most " + fieldError.Param() + " characters long"
	case "eqfield":
		return "This field must be equal to " + fieldError.Param()
	case "unique":
		return "This field must be unique"
	case "password":
		return "Password minimum length is 8"
	case "oneof":
		return "This field must be one of " + fieldError.Param()
	default:
		return "Invalid value"
	}
}

// ValidationError handles validation errors
func ValidationError(ctx echo.Context, err error) error {

	// Get detailed validation error messages
	validationErrors := err.(validator.ValidationErrors)
	errors := make(map[string]string)

	for _, fieldError := range validationErrors {
		errors[fieldError.Field()] = getErrorMessage(fieldError)
	}

	return ctx.JSON(http.StatusBadRequest, map[string]any{
		"error":   "Validation failed",
		"details": errors,
	})
}