package controllers

import (
	"log"
	"net/http"

	"kopelko-dating-app-backend/dto"
	m "kopelko-dating-app-backend/middlewares"
	"kopelko-dating-app-backend/services"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	userService services.AuthService
}

func NewAuthController(userService services.AuthService) *AuthController {
	return &AuthController{userService}
}

func (c *AuthController) RegisterUser(ctx echo.Context) error {
	log.Println("Attempting to register user")
	var req dto.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		log.Printf("Could not bind request: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := ctx.Validate(&req); err != nil {
		errors := m.GetValidationError(ctx, err)
		log.Printf("Validation failed: %v", errors)
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": errors,
		})
	}

	user, err := c.userService.RegisterUser(&req)
	if err != nil {
		log.Printf("Registration failed: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Printf("User registered successfully: %s", user.MaskEmail())
	return ctx.JSON(http.StatusCreated, map[string]any{
		"message": "User registered successfully",
		"user": map[string]any{
			"id":    "******",
			"email": user.MaskEmail(),
		},
	})
}

func (c *AuthController) LoginUser(ctx echo.Context) error {
	log.Println("Attempting to login user")
	var req dto.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		log.Printf("Could not bind request: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := ctx.Validate(&req); err != nil {
		errors := m.GetValidationError(ctx, err)
		log.Printf("Validation failed: %v", errors)
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": errors,
		})
	}

	user, err := c.userService.LoginUser(&req)
	if err != nil {
		log.Printf("Login failed: %v", err)
		return ctx.JSON(http.StatusUnauthorized, map[string]any{
			"error":   "Validation failed",
			"details": "Invalid email or password",
		})
	}

	ctx.Set("token", user.Token)
	log.Println("User logged in successfully")
	return ctx.JSON(http.StatusOK, map[string]any{
		"message": "User logged in successfully",
		"user": map[string]any{
			"id":    "******",
			"email": user.MaskEmail(),
		},
		"token": user.Token,
	})
}
