package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestCustomValidator_Validate(t *testing.T) {
	cv := NewValidator()

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			input: TestStruct{
				Email: "john.doe@example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: TestStruct{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cv.Validate(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
