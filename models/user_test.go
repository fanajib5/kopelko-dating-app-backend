package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MaskLocalPart(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "local part less than or equal to 3 characters",
			input:    "abc",
			expected: "abc",
		},
		{
			name:     "local part more than 3 characters",
			input:    "abcdef",
			expected: "abc***",
		},
		{
			name:     "local part exactly 3 characters",
			input:    "xyz",
			expected: "xyz",
		},
		{
			name:     "local part with special characters",
			input:    "a.b-c",
			expected: "a.b**",
		},
		{
			name:     "local part empty",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := maskLocalPart(tc.input)
			assert.Equal(t, tc.expected, output)
		})
	}
}

func Test_MaskEmail(t *testing.T) {
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
			expected: "a.b**@*****",
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
