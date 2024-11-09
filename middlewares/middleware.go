package utils

import (
	"kopelko-dating-app-backend/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware is a middleware function that validates the JWT token
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	jwtKey := utils.GetJWTKey()
	return func(ctx echo.Context) error {
		ctx.Logger().Print("Authenticating user...")

		// Retrieve the token from the Authorization header
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header missing"})
		}

		// Extract the token, expecting the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization format"})
		}
		tokenStr := parts[1]

		// Parse and validate the JWT token
		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		}

		// Token is valid; attach user information to the context
		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)

		// Proceed to the next handler
		return next(ctx)
	}
}