package middleware

import (
	"strings"

	pkgHttp "kopelko-dating-app-backend/internal/platform/http"
	"kopelko-dating-app-backend/internal/platform/token"

	"github.com/labstack/echo/v4"
)

const (
	ContextUserIDKey = "user_id"
	ContextEmailKey  = "user_email"
)

func AuthMiddleware(tokenSvc token.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return pkgHttp.Unauthorized(c, "Authorization header is required")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return pkgHttp.Unauthorized(c, "Invalid authorization format, must be Bearer token")
			}

			claims, err := tokenSvc.ValidateToken(parts[1])
			if err != nil {
				return pkgHttp.Unauthorized(c, "Invalid or expired token")
			}

			c.Set(ContextUserIDKey, claims.UserID)
			c.Set(ContextEmailKey, claims.Email)

			return next(c)
		}
	}
}

func GetCurrentUserID(c echo.Context) (uint, bool) {
	val := c.Get(ContextUserIDKey)
	if val == nil {
		return 0, false
	}
	id, ok := val.(uint)
	return id, ok
}
