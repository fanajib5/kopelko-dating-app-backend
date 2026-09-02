package middleware

import (
	"context"
	"time"

	pkgLogger "kopelko-dating-app-backend/internal/platform/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const HeaderXRequestID = "X-Request-ID"

func RequestIDAndLoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header.Get(HeaderXRequestID)
			if reqID == "" {
				reqID = uuid.New().String()
			}
			c.Response().Header().Set(HeaderXRequestID, reqID)

			ctx := context.WithValue(c.Request().Context(), pkgLogger.RequestIDKey, reqID)
			c.SetRequest(c.Request().WithContext(ctx))

			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			logger := pkgLogger.FromContext(ctx)
			logger.Info("HTTP request",
				"method", c.Request().Method,
				"uri", c.Request().RequestURI,
				"status", c.Response().Status,
				"latency_ms", latency.Milliseconds(),
				"ip", c.RealIP(),
			)

			return err
		}
	}
}
