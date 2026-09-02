package middleware

import (
	"net/http"
	"time"

	pkgHttp "kopelko-dating-app-backend/internal/platform/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// SecurityHeadersMiddleware adds standard web security headers
func SecurityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			res.Header().Set("X-Content-Type-Options", "nosniff")
			res.Header().Set("X-Frame-Options", "DENY")
			res.Header().Set("X-XSS-Protection", "1; mode=block")
			res.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			res.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			return next(c)
		}
	}
}

// NewAuthRateLimiter creates an in-memory rate limiter for sensitive authentication endpoints
func NewAuthRateLimiter(maxRequests int, window time.Duration) echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(float64(maxRequests) / window.Seconds()),
				Burst:     maxRequests,
				ExpiresIn: window * 2,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return pkgHttp.Error(context, http.StatusTooManyRequests, "Too many requests. Please try again later.", nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return pkgHttp.Error(context, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.", nil)
		},
	}
	return middleware.RateLimiterWithConfig(config)
}
