package controllers

import (
	"net/http"

	"kopelko-dating-app-backend/dto"
	"kopelko-dating-app-backend/services"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	userService *services.AuthService
}

func NewAuthController(userService *services.AuthService) *AuthController {
	return &AuthController{userService}
}

func (c *AuthController) RegisterUser(ctx echo.Context) error {
	var req dto.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	user, err := c.userService.RegisterUser(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
