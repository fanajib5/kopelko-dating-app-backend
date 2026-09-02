package http

import (
	"net/http"

	"kopelko-dating-app-backend/internal/plugins/subscription/domain"
	pkgHttp "kopelko-dating-app-backend/internal/platform/http"
	"kopelko-dating-app-backend/internal/platform/middleware"

	"github.com/labstack/echo/v4"
)

type SubscriptionHandler struct {
	svc domain.SubscriptionService
}

func NewSubscriptionHandler(svc domain.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

type SubscribeRequest struct {
	FeatureName string `json:"feature_name" validate:"required" example:"no_swipe_quota"`
}

// Subscribe godoc
// @Summary Subscribe to premium feature
// @Description Purchase premium feature such as 'no_swipe_quota' or 'verified_label'
// @Tags Subscriptions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body SubscribeRequest true "Subscribe payload"
// @Success 200 {object} pkgHttp.APIResponse{data=domain.Subscription}
// @Failure 400 {object} pkgHttp.APIResponse
// @Failure 401 {object} pkgHttp.APIResponse
// @Router /api/users/subscriptions [post]
func (h *SubscriptionHandler) Subscribe(c echo.Context) error {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		return pkgHttp.Unauthorized(c, "Unauthorized")
	}

	var req SubscribeRequest
	if err := c.Bind(&req); err != nil {
		return pkgHttp.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return pkgHttp.BadRequest(c, "Validation failed", err.Error())
	}

	sub, err := h.svc.Subscribe(c.Request().Context(), userID, req.FeatureName)
	if err != nil {
		return pkgHttp.BadRequest(c, err.Error(), nil)
	}

	return pkgHttp.Success(c, http.StatusOK, "Subscription activated successfully", sub)
}

func (h *SubscriptionHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/subscriptions", h.Subscribe)
}
