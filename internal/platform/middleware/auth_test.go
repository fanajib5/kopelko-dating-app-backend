package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kopelko-dating-app-backend/internal/platform/middleware"
	"kopelko-dating-app-backend/internal/platform/token"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	jwtSvc := token.NewJWTService("test_secret_key_1234567890123456")
	validToken, err := jwtSvc.GenerateToken(42, "tester@example.com")
	assert.NoError(t, err)

	e := echo.New()

	handler := func(c echo.Context) error {
		userID, ok := middleware.GetCurrentUserID(c)
		assert.True(t, ok)
		assert.Equal(t, uint(42), userID)
		return c.String(http.StatusOK, "authorized")
	}

	mw := middleware.AuthMiddleware(jwtSvc)

	t.Run("Valid Bearer Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := mw(handler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := mw(handler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Invalid Token Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Basic "+validToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := mw(handler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
