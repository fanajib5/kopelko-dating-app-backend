package utils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_utils_InitDB(t *testing.T) {
	// Create a mock database connection
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	// Use the mock database connection with GORM
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	require.NoError(t, err)
	assert.NotNil(t, db)

	t.Cleanup(func() {
		mockDB.Close()
	})
}
