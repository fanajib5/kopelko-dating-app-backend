package config

import (
	"os"
	"testing"

	"kopelko-dating-app-backend/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_NewValidator(t *testing.T) {
	expected := utils.NewValidator()

	output := NewValidator()
	assert.NotNil(t, output)
	assert.IsType(t, expected, output)
}

func Test_LoadEnv(t *testing.T) {
	err := os.WriteFile(".env", []byte("API_PORT=8080"), 0644)
	require.NoError(t, err)
	loadEnv()

	testEnv := os.Getenv("API_PORT")
	assert.Equal(t, "8080", testEnv)

	t.Cleanup(func() {
		os.Remove(".env")
	})
}

func Test_LoadAPIPort(t *testing.T) {
	// Test when API_PORT is set in the environment
	err := os.Setenv("API_PORT", "8080")
	require.NoError(t, err)

	config := &Config{}
	config.LoadAPIPort()
	assert.Equal(t, "8080", config.APIPort)

	// Test when API_PORT is not set in the environment
	err = os.Unsetenv("API_PORT")
	require.NoError(t, err)

	config.LoadAPIPort()
	assert.Empty(t, config.APIPort)

	t.Cleanup(func() {
		os.Unsetenv("API_PORT")
	})
}

func Test_initializeControllers(t *testing.T) {
	// Create a mock database connection
	db, _, err := sqlmock.New()
	require.NoError(t, err)

	// Use the mock database connection with GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	config := &Config{DB: gormDB, LimitSwipe: 10}
	config.initializeControllers()

	assert.NotNil(t, config.Controllers.Profile)
	assert.NotNil(t, config.Controllers.Auth)
	assert.NotNil(t, config.Controllers.Subscribe)
	assert.NotNil(t, config.Controllers.Swipe)

	t.Cleanup(func() {
		db.Close()
	})
}

func Test_LoadLimitSwipe(t *testing.T) {
	// Test when LIMIT_SWIPE is set and valid
	err := os.Setenv("LIMIT_SWIPE", "15")
	require.NoError(t, err)

	config := &Config{}
	config.LoadLimitSwipe()
	assert.Equal(t, 15, config.LimitSwipe)

	// Test when LIMIT_SWIPE is set but invalid
	err = os.Setenv("LIMIT_SWIPE", "invalid")
	require.NoError(t, err)

	config.LoadLimitSwipe()
	assert.Equal(t, defaultLimitSwipe, config.LimitSwipe)

	// Test when LIMIT_SWIPE is not set
	err = os.Unsetenv("LIMIT_SWIPE")
	require.NoError(t, err)

	config.LoadLimitSwipe()
	assert.Equal(t, defaultLimitSwipe, config.LimitSwipe)

	t.Cleanup(func() {
		os.Unsetenv("LIMIT_SWIPE")
	})
}

func Test_LoadJWTKey(t *testing.T) {
	// Mock the LoadJWTKey function
	expectedKey := []byte("mocked_jwt_key")

	err := os.Setenv("JWT_KEY", "mocked_jwt_key")
	require.NoError(t, err)

	config := &Config{}
	config.LoadJWTKey()
	assert.Equal(t, expectedKey, config.JWTKey)
}
