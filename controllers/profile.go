package controllers

import (
	"net/http"

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
	id := ctx.Get("user_id").(uint)

	profile, err := c.profileService.GetProfileByID(id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
	}

	return ctx.JSON(http.StatusOK, profile)
}

func (c *ProfileController) RandomProfile(ctx echo.Context) error {
	profile, err := c.profileService.GetRandomProfile()
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
	}

	return ctx.JSON(http.StatusOK, profile)
}
