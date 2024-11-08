package controllers

import (
	"errors"
	"net/http"
	"strconv"

	service "kopelko-dating-app-backend/services"

	"github.com/labstack/echo/v4"
)

type SwipeController struct {
	swipeService service.SwipeService
}

// NewSwipeController creates a new SwipeController
func NewSwipeController(swipeService service.SwipeService) *SwipeController {
	return &SwipeController{swipeService: swipeService}
}

// SwipeHandler processes swipe requests
func (sc *SwipeController) SwipeHandler(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uint)
	isPremium, ok := ctx.Get("is_premium").(bool)
	if !ok {
		ctx.Logger().Error(errors.New("Could not get user premium status"))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get user premium status"})
	}

	targetUserID, err := strconv.Atoi(ctx.Param("target_user_id"))
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target user ID"})
	}

	swipeType := ctx.QueryParam("type")
	if swipeType != "left" && swipeType != "right" {
		ctx.Logger().Error(errors.New("Invalid swipe type, must be 'left' or 'right'"))
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid swipe type"})
	}

	err = sc.swipeService.ProcessSwipe(userID, targetUserID, swipeType, isPremium)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Swipe successful"})
}
