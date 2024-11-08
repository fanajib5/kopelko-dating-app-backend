package controllers

import (
	"net/http"

	"kopelko-dating-app-backend/dto"
	service "kopelko-dating-app-backend/services"
	util "kopelko-dating-app-backend/utils"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	userService service.AuthService
}

func NewAuthController(userService service.AuthService) *AuthController {
	return &AuthController{userService}
}

func (c *AuthController) RegisterUser(ctx echo.Context) error {
	var req dto.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.Logger().Errorf("Could not bind request: %w", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := ctx.Validate(&req); err != nil {
		errors := util.ValidationError(ctx, err)
		ctx.Logger().Errorf("Validation failed: %w", errors)
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": errors,
		})
	}

	user, err := c.userService.RegisterUser(&req)
	if err != nil {
		ctx.Logger().Errorf("Registration failed: %w", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	ctx.Logger().Printf("User registered successfully: %s", user.MaskEmail())
	return ctx.JSON(http.StatusCreated, map[string]any{
		"message": "User registered successfully",
		"user": map[string]any{
			"id":    "******",
			"email": user.MaskEmail(),
		},
	})
}

func (c *AuthController) LoginUser(ctx echo.Context) error {
	var req dto.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.Logger().Errorf("Could not bind request: %w", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := ctx.Validate(&req); err != nil {
		errors := util.ValidationError(ctx, err)
		ctx.Logger().Errorf("Validation failed: %w", errors)
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": errors,
		})
	}

	user, err := c.userService.LoginUser(&req)
	if err != nil {
		ctx.Logger().Errorf("Login failed: %w", err)
		return ctx.JSON(http.StatusUnauthorized, map[string]any{
			"error":   "Validation failed",
			"details": "Invalid email or password",
		})
	}

	ctx.Set("token", user.Token)
	ctx.Logger().Print("User logged in successfully")
	return ctx.JSON(http.StatusOK, map[string]any{
		"message": "User logged in successfully",
		"user": map[string]any{
			"id":    "******",
			"email": user.MaskEmail(),
		},
		"token": user.Token,
	})
}
