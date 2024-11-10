package controllers

import (
	"log"
	"net/http"
	"strconv"

	m "kopelko-dating-app-backend/middlewares"
	"kopelko-dating-app-backend/models"
	"kopelko-dating-app-backend/services"

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
	log.Println("Attempting to swipe an user")

	userID := m.GetUserIDFromContext(ctx)

	targetUserID, err := strconv.Atoi(ctx.Param("target_user_id"))
	if err != nil {
		log.Printf("Invalid target user ID: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target user ID"})
	}

	swipeType := ctx.QueryParam("type")
	if swipeType != models.SwipeTypePass && swipeType != models.SwipeTypeLike {
		errMsg := "Invalid swipe type, must be 'pass' or 'like'"
		log.Println(errMsg)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
	}

	err = sc.swipeService.SwipeProfile(userID, targetUserID, swipeType)
	if err != nil {
		ec, errMsg := m.ParseErrorCodeAndMessage(err)
		log.Println("Failed to swipe profile:", errMsg)
		return ctx.JSON(ec, map[string]string{"error": errMsg})
	}

	log.Println("Swipe an user successful")
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Swipe successful"})
}
