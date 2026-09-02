package http

import (
	"net/http"
	"strconv"

	"kopelko-dating-app-backend/internal/plugins/swipe/domain"
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
	SwipeType string `json:"swipe_type" validate:"required,oneof=like pass" example:"like"`
}

// Swipe godoc
// @Summary Swipe a candidate user profile
// @Description Record swipe decision (like or pass). If mutual like occurs, returns match metadata. Limited to 10/day unless subscribed.
// @Tags Swipes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param target_user_id path int true "Target User ID to swipe on"
// @Param request body SwipeRequest true "Swipe payload"
// @Success 200 {object} pkgHttp.APIResponse{data=domain.SwipeResponse}
// @Failure 400 {object} pkgHttp.APIResponse
// @Failure 401 {object} pkgHttp.APIResponse
// @Router /api/users/swipes/{target_user_id} [post]
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

// GetMatches godoc
// @Summary Get all matches
// @Description Retrieve a list of mutual matches along with candidate profile metadata
// @Tags Swipes
// @Security BearerAuth
// @Produce json
// @Success 200 {object} pkgHttp.APIResponse{data=[]domain.MatchDetail}
// @Failure 401 {object} pkgHttp.APIResponse
// @Router /api/users/matches [get]
func (h *SwipeHandler) GetMatches(c echo.Context) error {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		return pkgHttp.Unauthorized(c, "Unauthorized")
	}

	matches, err := h.svc.GetMatches(c.Request().Context(), userID)
	if err != nil {
		return pkgHttp.InternalServerError(c, "Failed to retrieve matches", err.Error())
	}

	if matches == nil {
		matches = []domain.MatchDetail{}
	}

	return pkgHttp.Success(c, http.StatusOK, "Matches retrieved successfully", matches)
}

func (h *SwipeHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/swipes/:target_user_id", h.Swipe)
	g.GET("/matches", h.GetMatches)
}
