package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()
	mockToken := func() string {
		token, err := utils.GenerateJWT(models.User{ID: 1, Email: "test@xample.com"})
		if err != nil {
			t.Fatal(err)
		}
		return token
	}()

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Missing Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Authorization header missing"}`,
		},
		{
			name:           "Invalid Authorization Format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid authorization format"}`,
		},
		{
			name:           "Invalid Token",
			authHeader:     "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid token"}`,
		},
		{
			name:           "Valid Token",
			authHeader:     "Bearer " + mockToken,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Mock the next handler
			next := func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{"message": "success"})
			}

			handler := AuthMiddleware(next)
			err := handler(ctx)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
