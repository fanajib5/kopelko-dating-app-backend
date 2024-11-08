package controllers

import (
	"net/http"

	"kopelko-dating-app-backend/dto"
	"kopelko-dating-app-backend/services"
	"kopelko-dating-app-backend/utils"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	userService services.AuthService
}

func NewAuthController(userService services.AuthService) *AuthController {
	return &AuthController{userService}
}

func (c *AuthController) RegisterUser(ctx echo.Context) error {
	var req dto.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ctx.Validate(&req); err != nil {
		ctx.Logger().Error(err)
		return utils.ValidationError(ctx, err)
	}

	user, err := c.userService.RegisterUser(&req)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	ctx.Logger().Info("User registered successfully")

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
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ctx.Validate(&req); err != nil {
		ctx.Logger().Error(err)
		return utils.ValidationError(ctx, err)
	}

	user, err := c.userService.LoginUser(&req)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusUnauthorized, map[string]any{
			"error":   "Validation failed",
			"details": "Invalid email or password",
		})
	}

	ctx.Set("token", user.Token)
	ctx.Logger().Info("User logged in successfully")

	return ctx.JSON(http.StatusOK, map[string]any{
		"message": "User logged in successfully",
		"user": map[string]any{
			"id":    user.ID,
			"email": user.MaskEmail(),
		},
		"token": user.Token, // Assuming the user service returns a token
	})
}
