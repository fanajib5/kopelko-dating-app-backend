package subscription

import (
	"context"

	"kopelko-dating-app-backend/internal/core/plugin"
	subHttp "kopelko-dating-app-backend/internal/plugins/subscription/delivery/http"
	"kopelko-dating-app-backend/internal/plugins/subscription/domain"
	subRepo "kopelko-dating-app-backend/internal/plugins/subscription/repository"
	subUsecase "kopelko-dating-app-backend/internal/plugins/subscription/usecase"

	"github.com/labstack/echo/v4"
)

type SubscriptionPlugin struct {
	repo    domain.SubscriptionRepository
	service domain.SubscriptionService
	handler *subHttp.SubscriptionHandler
}

func NewPlugin() plugin.Plugin {
	return &SubscriptionPlugin{}
}

func (p *SubscriptionPlugin) Name() string {
	return "subscription"
}

func (p *SubscriptionPlugin) Version() string {
	return "1.0.0"
}

func (p *SubscriptionPlugin) Init(appCtx *plugin.AppContext) error {
	p.repo = subRepo.NewSubscriptionRepository(appCtx.DBPool)
	p.service = subUsecase.NewSubscriptionUsecase(p.repo)
	p.handler = subHttp.NewSubscriptionHandler(p.service)

	// Register Filter Hook: "swipe.check_quota"
	// Payload passed is map[string]any{"user_id": uint, "has_unlimited": bool}
	if appCtx.Hooks != nil {
		appCtx.Hooks.AddFilter("swipe.check_quota", 10, func(ctx context.Context, data any) (any, error) {
			m, ok := data.(map[string]any)
			if !ok {
				return data, nil
			}
			userID, _ := m["user_id"].(uint)
			hasNoQuota, err := p.service.HasActiveFeature(ctx, userID, "no_swipe_quota")
			if err == nil && hasNoQuota {
				m["has_unlimited"] = true
			}
			return m, nil
		})

		// Register Filter Hook: "profile.decorate"
		// Payload passed is map[string]any{"user_id": uint, "is_verified": bool, "is_premium": bool}
		appCtx.Hooks.AddFilter("profile.decorate", 10, func(ctx context.Context, data any) (any, error) {
			m, ok := data.(map[string]any)
			if !ok {
				return data, nil
			}
			userID, _ := m["user_id"].(uint)
			isVerified, _ := p.service.HasActiveFeature(ctx, userID, "verified_label")
			if isVerified {
				m["is_verified"] = true
			}
			isUnlimited, _ := p.service.HasActiveFeature(ctx, userID, "no_swipe_quota")
			if isUnlimited || isVerified {
				m["is_premium"] = true
			}
			return m, nil
		})
	}

	return nil
}

func (p *SubscriptionPlugin) RegisterRoutes(apiGroup *echo.Group, authMiddleware echo.MiddlewareFunc) {
	// Subscriptions are mounted under protected users route
	userGroup := apiGroup.Group("/users", authMiddleware)
	p.handler.RegisterRoutes(userGroup)
}

func (p *SubscriptionPlugin) Service() domain.SubscriptionService {
	return p.service
}
