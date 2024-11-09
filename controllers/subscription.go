package controllers

import (
	"net/http"
	"strconv"

	service "kopelko-dating-app-backend/services"

	"github.com/labstack/echo/v4"
)

type SubscriptionController struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionController(service service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{subscriptionService: service}
}

func (sc *SubscriptionController) SubscribeHandler(ctx echo.Context) error {
	userID, err := strconv.Atoi(ctx.Get("user_id").(string))
	if err != nil {
		ctx.Logger().Errorf("Could not get user ID: %w", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Could not get user ID"})
	}

	featureID, err := strconv.Atoi(ctx.QueryParam("feature_id"))
	if err != nil {
		ctx.Logger().Errorf("Invalid feature ID: %w", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid feature ID"})
	}

	err = sc.subscriptionService.SubscribeUser(userID, featureID)
	if err != nil {
		ctx.Logger().Errorf("Subscription failed: %w", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	ctx.Logger().Printf("User %d subscribed to feature %d succesfully", userID, featureID)
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Subscription successful"})
}
