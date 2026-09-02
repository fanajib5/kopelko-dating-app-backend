package middleware

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			for _, valErr := range valErrs {
				return fmt.Errorf("field '%s' failed on rule '%s'", valErr.Field(), valErr.Tag())
			}
		}
		return err
	}
	return nil
}
