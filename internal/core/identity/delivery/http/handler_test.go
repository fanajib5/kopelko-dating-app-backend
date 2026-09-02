package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	deliveryHttp "kopelko-dating-app-backend/internal/core/identity/delivery/http"
	"kopelko-dating-app-backend/internal/core/identity/domain"
	pkgHttp "kopelko-dating-app-backend/internal/platform/http"
	"kopelko-dating-app-backend/internal/platform/middleware"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIdentityService struct {
	mock.Mock
}

func (m *MockIdentityService) Register(ctx context.Context, email, password, name string, age int, gender, location string, interests, photos []string) (*domain.User, string, error) {
	args := m.Called(ctx, email, password, name, age, gender, location, interests, photos)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*domain.User), args.String(1), args.Error(2)
}

func (m *MockIdentityService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*domain.User), args.String(1), args.Error(2)
}

func TestIdentityHandler_Register_HTTP(t *testing.T) {
	e := echo.New()
	e.Validator = middleware.NewValidator()
	svc := new(MockIdentityService)
	handler := deliveryHttp.NewIdentityHandler(svc)

	t.Run("Success 201", func(t *testing.T) {
		reqBody := deliveryHttp.RegisterRequest{
			Email:    "new@example.com",
			Password: "password123",
			Name:     "New User",
			Age:      24,
			Gender:   "male",
			Location: "Jakarta",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		user := &domain.User{ID: 1, Email: "new@example.com"}
		svc.On("Register", mock.Anything, "new@example.com", "password123", "New User", 24, "male", "Jakarta", []string(nil), []string(nil)).
			Return(user, "jwt-test-token", nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp pkgHttp.APIResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("Validation Error 400", func(t *testing.T) {
		reqBody := deliveryHttp.RegisterRequest{
			Email:    "invalid-email",
			Password: "123", // too short
		}
		jsonBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestIdentityHandler_Login_HTTP(t *testing.T) {
	e := echo.New()
	e.Validator = middleware.NewValidator()
	svc := new(MockIdentityService)
	handler := deliveryHttp.NewIdentityHandler(svc)

	t.Run("Success 200", func(t *testing.T) {
		reqBody := deliveryHttp.LoginRequest{
			Email:    "user@example.com",
			Password: "secretpassword",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		user := &domain.User{ID: 1, Email: "user@example.com"}
		svc.On("Login", mock.Anything, "user@example.com", "secretpassword").
			Return(user, "jwt-login-token", nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Unauthorized 401", func(t *testing.T) {
		reqBody := deliveryHttp.LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		svc.On("Login", mock.Anything, "user@example.com", "wrongpassword").
			Return(nil, "", errors.New("invalid email or password")).Once()

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
