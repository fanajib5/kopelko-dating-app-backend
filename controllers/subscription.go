package controllers

import (
	"net/http"
	"strconv"

	"kopelko-dating-app-backend/services"

	"github.com/labstack/echo/v4"
)

type SubscriptionController struct {
	subscriptionService services.SubscriptionService
}

func NewSubscriptionController(service services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{subscriptionService: service}
}

func (sc *SubscriptionController) SubscribeHandler(ctx echo.Context) error {
	userIDjwt := ctx.Get("user_id").(uint)
	featureID, err := strconv.Atoi(ctx.QueryParam("feature_id"))
	if err != nil {
		ctx.Logger().Errorf("Invalid feature ID: %w", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid feature ID"})
	}

	err = sc.subscriptionService.SubscribeUser(userIDjwt, featureID)
	if err != nil {
		ctx.Logger().Errorf("Subscription failed: %w", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	ctx.Logger().Printf("User %d subscribed to feature %d succesfully", userIDjwt, featureID)
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Subscription successful"})
}
