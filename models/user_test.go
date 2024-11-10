package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kopelko-dating-app-backend/utils"
)

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "Empty email",
			email:    "",
			expected: "",
		},
		{
			name:     "Valid email",
			email:    "test@example.com",
			expected: utils.MaskEmail("test@example.com"),
		},
		{
			name:     "Another valid email",
			email:    "user@domain.com",
			expected: utils.MaskEmail("user@domain.com"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{Email: tt.email}
			result := user.MaskEmail()
			assert.Equal(t, tt.expected, result)
		})
	}
}
