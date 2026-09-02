package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kopelko-dating-app-backend/internal/platform/middleware"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}

	mw := middleware.SecurityHeadersMiddleware()
	err := mw(handler)(c)

	assert.NoError(t, err)
	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", rec.Header().Get("X-XSS-Protection"))
	assert.Contains(t, rec.Header().Get("Strict-Transport-Security"), "max-age=31536000")
}

func TestNewAuthRateLimiter(t *testing.T) {
	e := echo.New()
	limiter := middleware.NewAuthRateLimiter(2, 1*time.Minute)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}

	// Request 1: OK
	req1 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req1.RemoteAddr = "192.168.1.1:1234"
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req1, rec1)
	err1 := limiter(handler)(c1)
	assert.NoError(t, err1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// Request 2: OK
	req2 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req2.RemoteAddr = "192.168.1.1:1234"
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	err2 := limiter(handler)(c2)
	assert.NoError(t, err2)
	assert.Equal(t, http.StatusOK, rec2.Code)

	// Request 3: Rate limited (429)
	req3 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req3.RemoteAddr = "192.168.1.1:1234"
	rec3 := httptest.NewRecorder()
	c3 := e.NewContext(req3, rec3)
	err3 := limiter(handler)(c3)
	assert.NoError(t, err3)
	assert.Equal(t, http.StatusTooManyRequests, rec3.Code)
}
