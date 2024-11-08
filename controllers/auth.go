package controllers

import (
	"net/http"

	"kopelko-dating-app-backend/dto"
	"kopelko-dating-app-backend/services"
	"kopelko-dating-app-backend/utils"

	"github.com/go-playground/validator/v10"
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
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ctx.Validate(&req); err != nil {
		ctx.Logger().Error(err)

		// Get detailed validation error messages
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)

		for _, fieldError := range validationErrors {
			errors[fieldError.Field()] = utils.GetErrorMessage(fieldError)
		}

		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": errors,
		})
	}

	user, err := c.userService.RegisterUser(&req)
	if err != nil {
		ctx.Logger().Error(err)
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
