package controllers

import (
	"errors"
	"log"
	"net/http"

	m "kopelko-dating-app-backend/middlewares"
	"kopelko-dating-app-backend/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProfileController struct {
	profileService services.ProfileService
}

func NewProfileController(profileService services.ProfileService) *ProfileController {
	return &ProfileController{profileService: profileService}
}

func (c *ProfileController) ViewMyProfile(ctx echo.Context) error {
	log.Println("Attempting to view my profile")

	// Retrieve the user ID from the Echo context
	id := m.GetUserIDFromContext(ctx)

	profile, err := c.profileService.GetProfileByID(id)
	if err != nil {
		log.Println("Failed to get profile:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Println("Profile retrieved successfully")
	return ctx.JSON(http.StatusOK, profile)
}

func (c *ProfileController) RandomProfiles(ctx echo.Context) error {
	log.Println("Attempting to get random profiles")
	viewerID := m.GetUserIDFromContext(ctx)

	profiles, err := c.profileService.GetRandomProfiles(viewerID)
	if err != nil {
		ec, errMsg := m.ParseErrorCodeAndMessage(err)
		log.Println("Failed to get random profiles:", errMsg)
		return ctx.JSON(ec, map[string]string{"error": errMsg})
	}

	log.Println("Random profiles retrieved successfully")
	return ctx.JSON(http.StatusOK, profiles)
}
