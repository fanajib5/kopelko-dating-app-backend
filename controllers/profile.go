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

func (c *ProfileController) ViewProfile(ctx echo.Context) error {
	id := ctx.Param("id")
	profile, err := c.profileService.GetProfileByID(id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
	}

	return ctx.JSON(http.StatusOK, profile)
}
