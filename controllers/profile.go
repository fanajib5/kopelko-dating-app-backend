package controllers

import (
	"net/http"

	m "kopelko-dating-app-backend/middlewares"
	"kopelko-dating-app-backend/services"

	"github.com/labstack/echo/v4"
)

type ProfileController struct {
	profileService services.ProfileService
}

func NewProfileController(profileService services.ProfileService) *ProfileController {
	return &ProfileController{profileService: profileService}
}

func (c *ProfileController) ViewMyProfile(ctx echo.Context) error {
	// Retrieve the user ID from the Echo context
	id := m.GetUserIDFromContext(ctx)

	profile, err := c.profileService.GetProfileByID(id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
	}

	return ctx.JSON(http.StatusOK, profile)
}

func (c *ProfileController) RandomProfiles(ctx echo.Context) error {
	viewerID := m.GetUserIDFromContext(ctx)
	profiles, err := c.profileService.GetRandomProfiles(viewerID)
	if err != nil {
		ec, errMsg := m.ParseErrorCodeAndMessage(err)
		ctx.Logger().Error(errMsg)
		return ctx.JSON(ec, map[string]string{"error": errMsg})
	}

	return ctx.JSON(http.StatusOK, profiles)
}
