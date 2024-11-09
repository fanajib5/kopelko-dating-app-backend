package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_models_MaskEmail(t *testing.T) {
	testCases := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "valid email with local part more than 3 characters",
			email:    "test@example.com",
			expected: "tes*****@*****",
		},
		{
			name:     "valid email with local part exactly 3 characters",
			email:    "abc@example.com",
			expected: "abc@*****",
		},
		{
			name:     "valid email with local part less than 3 characters",
			email:    "ab@example.com",
			expected: "ab@*****",
		},
		{
			name:     "email without @ character",
			email:    "invalidemail",
			expected: "invalidemail",
		},
		{
			name:     "email with special characters in local part",
			email:    "a.b-c@example.com",
			expected: "a.b*****@*****",
		},
		{
			name:     "empty email",
			email:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := User{Email: tc.email}
			output := user.MaskEmail()
			assert.Equal(t, tc.expected, output)
		})
	}
}
