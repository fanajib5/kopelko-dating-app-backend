package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/services"
	"kopelko-dating-app-backend/utils"

	"github.com/labstack/echo/v4"
)

type SwipeController struct {
	swipeService services.SwipeService
}

// NewSwipeController creates a new SwipeController
func NewSwipeController(swipeService services.SwipeService) *SwipeController {
	return &SwipeController{swipeService: swipeService}
}

// SwipeHandler processes swipe requests
func (sc *SwipeController) SwipeHandler(ctx echo.Context) error {
	userID := utils.GetUserIDFromContext(ctx)

	targetUserID, err := strconv.Atoi(ctx.Param("target_user_id"))
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target user ID"})
	}

	swipeType := ctx.QueryParam("type")
	if swipeType != models.SwipeTypePass && swipeType != models.SwipeTypeLike {
		ctx.Logger().Error(errors.New("Invalid swipe type, must be 'pass' or 'like'"))
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid swipe type"})
	}

	err = sc.swipeService.SwipeProfile(userID, targetUserID, swipeType)
	if err != nil {
		ctx.Logger().Error(err)
		ec, errMsg := utils.ParseErrorCodeAndMessage(err)
		return ctx.JSON(ec, map[string]string{"error": errMsg})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Swipe successful"})
}
