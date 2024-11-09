package utils

import (
	"os"
	"strings"
	"testing"
	"time"

	"kopelko-dating-app-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_utils_LoadJWTKey(t *testing.T) {
	expectedKey := "mysecretkey"
	os.Setenv("JWT_KEY", expectedKey)

	actualKey := LoadJWTKey()
	assert.Equal(t, []byte(expectedKey), actualKey)
}

func Test_utils_GenerateJWT(t *testing.T) {
	LoadJWTKey()
	user := models.User{
		ID:    1,
		Email: "test@example.com",
	}

	tokenString, err := GenerateJWT(user)
	require.NoError(t, err)

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	require.NoError(t, err)

	assert.True(t, token.Valid)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, "kopelko-dating-app-backend", claims.Issuer)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func Test_utils_ParseJWT(t *testing.T) {
	LoadJWTKey()
	user := models.User{
		ID:    1,
		Email: "test@example.com",
	}

	tokenString, err := GenerateJWT(user)
	require.NoError(t, err)

	claims, err := ParseJWT(tokenString)
	require.NoError(t, err)

	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, "kopelko-dating-app-backend", claims.Issuer)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func Test_utils_ParseJWT_Failed(t *testing.T) {
	LoadJWTKey()
	expiredTime := time.Now().Add(-15 * time.Minute)
	claims := &Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "kopelko-dating-app-backend",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, err := token.SignedString(jwtKey)
	require.NoError(t, err)

	_, err = ParseJWT(expiredTokenString)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "failed to parsing JWT:"))
}
