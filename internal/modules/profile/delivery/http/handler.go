package http

import (
	"net/http"

	"kopelko-dating-app-backend/internal/modules/profile/domain"
	pkgHttp "kopelko-dating-app-backend/internal/platform/http"
	"kopelko-dating-app-backend/internal/platform/middleware"

	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	svc domain.ProfileService
}

func NewProfileHandler(svc domain.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) ViewMyProfile(c echo.Context) error {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		return pkgHttp.Unauthorized(c, "Unauthorized")
	}

	profile, err := h.svc.GetMyProfile(c.Request().Context(), userID)
	if err != nil {
		return pkgHttp.NotFound(c, "Profile not found")
	}

	return pkgHttp.Success(c, http.StatusOK, "Profile retrieved successfully", profile)
}

func (h *ProfileHandler) RandomProfiles(c echo.Context) error {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		return pkgHttp.Unauthorized(c, "Unauthorized")
	}

	profiles, err := h.svc.GetRandomProfiles(c.Request().Context(), userID)
	if err != nil {
		return pkgHttp.BadRequest(c, err.Error(), nil)
	}

	return pkgHttp.Success(c, http.StatusOK, "Profiles retrieved successfully", profiles)
}

func (h *ProfileHandler) RegisterRoutes(g *echo.Group) {
	profiles := g.Group("/profiles")
	profiles.GET("/me", h.ViewMyProfile)
	profiles.GET("/random", h.RandomProfiles)
}
