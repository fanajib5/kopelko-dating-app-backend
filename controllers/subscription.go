package controllers

import (
	"log"
	"net/http"
	"strconv"

	m "kopelko-dating-app-backend/middlewares"
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
	log.Println("Attempting to subscribe user to feature")

	userIDjwt := m.GetUserIDFromContext(ctx)
	featureID, err := strconv.Atoi(ctx.QueryParam("feature_id"))
	if err != nil {
		log.Printf("Invalid feature ID: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid feature ID"})
	}

	err = sc.subscriptionService.SubscribeUser(userIDjwt, featureID)
	if err != nil {
		log.Printf("Subscription failed: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Printf("User %d subscribed to feature %d succesfully", userIDjwt, featureID)
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Subscription successful"})
}
