package http

import (
	"net/http"
	"strconv"

	"kopelko-dating-app-backend/internal/modules/swipe/domain"
	pkgHttp "kopelko-dating-app-backend/internal/platform/http"
	"kopelko-dating-app-backend/internal/platform/middleware"

	"github.com/labstack/echo/v4"
)

type SwipeHandler struct {
	svc domain.SwipeService
}

func NewSwipeHandler(svc domain.SwipeService) *SwipeHandler {
	return &SwipeHandler{svc: svc}
}

type SwipeRequest struct {
	SwipeType string `json:"swipe_type" validate:"required,oneof=like pass"`
}

func (h *SwipeHandler) Swipe(c echo.Context) error {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		return pkgHttp.Unauthorized(c, "Unauthorized")
	}

	targetUserIDStr := c.Param("target_user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 64)
	if err != nil {
		return pkgHttp.BadRequest(c, "Invalid target_user_id parameter", err.Error())
	}

	var req SwipeRequest
	if err := c.Bind(&req); err != nil {
		return pkgHttp.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return pkgHttp.BadRequest(c, "Validation failed", err.Error())
	}

	res, err := h.svc.SwipeUser(c.Request().Context(), userID, uint(targetUserID), domain.SwipeType(req.SwipeType))
	if err != nil {
		return pkgHttp.BadRequest(c, err.Error(), nil)
	}

	return pkgHttp.Success(c, http.StatusOK, "Swipe processed successfully", res)
}

func (h *SwipeHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/swipes/:target_user_id", h.Swipe)
}
