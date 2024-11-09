package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(nil, nil)

	validate := validator.New()

	type TestStruct struct {
		Email string `validate:"required,email"`
	}

	testCases := []struct {
		name          string
		input         TestStruct
		expectedError map[string]string
	}{
		{
			name:  "Missing required field",
			input: TestStruct{},
			expectedError: map[string]string{
				"Email": "This field is required",
			},
		},
		{
			name:  "Invalid email format",
			input: TestStruct{Email: "invalid-email"},
			expectedError: map[string]string{
				"Email": "Invalid email address",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.input)
			if err != nil {
				validationErrors := ValidationError(ctx, err)
				assert.Equal(t, tc.expectedError, validationErrors)
			} else {
				t.Errorf("Expected validation error but got none")
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Set user_id in context
	expectedUserID := uint(123)
	ctx.Set("user_id", expectedUserID)

	// Retrieve user_id using GetUserIDFromContext
	userID := GetUserIDFromContext(ctx)

	// Assert that the retrieved user_id matches the expected value
	assert.Equal(t, expectedUserID, userID)
}

func TestParseErrorCodeAndMessage(t *testing.T) {
	testCases := []struct {
		name            string
		inputError      error
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "Echo HTTPError",
			inputError:      echo.NewHTTPError(http.StatusNotFound, "Resource not found"),
			expectedCode:    http.StatusNotFound,
			expectedMessage: "Resource not found",
		},
		{
			name:            "Generic error",
			inputError:      assert.AnError,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: assert.AnError.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code, message := ParseErrorCodeAndMessage(tc.inputError)
			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedMessage, message)
		})
	}
}
