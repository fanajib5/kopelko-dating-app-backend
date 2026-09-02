package profile

import (
	"context"
	"fmt"

	identityDomain "kopelko-dating-app-backend/internal/core/identity/domain"
	"kopelko-dating-app-backend/internal/core/plugin"
	profileHttp "kopelko-dating-app-backend/internal/plugins/profile/delivery/http"
	"kopelko-dating-app-backend/internal/plugins/profile/domain"
	profileRepo "kopelko-dating-app-backend/internal/plugins/profile/repository"
	profileUsecase "kopelko-dating-app-backend/internal/plugins/profile/usecase"

	"github.com/labstack/echo/v4"
)

type ProfilePlugin struct {
	repo    domain.ProfileRepository
	service domain.ProfileService
	handler *profileHttp.ProfileHandler
}

func NewPlugin() plugin.Plugin {
	return &ProfilePlugin{}
}

func (p *ProfilePlugin) Name() string {
	return "profile"
}

func (p *ProfilePlugin) Version() string {
	return "1.0.0"
}

func (p *ProfilePlugin) Init(appCtx *plugin.AppContext) error {
	p.repo = profileRepo.NewProfileRepository(appCtx.DBPool)
	p.service = profileUsecase.NewProfileUsecase(p.repo, appCtx.Hooks, appCtx.Config.LimitSwipe)
	p.handler = profileHttp.NewProfileHandler(p.service)

	// Action Hook listener: create profile on "user.registered"
	if appCtx.Hooks != nil {
		appCtx.Hooks.AddAction("user.registered", 10, func(ctx context.Context, payload any) error {
			regPayload, ok := payload.(*identityDomain.UserRegisteredPayload)
			if !ok {
				return nil
			}

			profile := &domain.Profile{
				UserID:    regPayload.User.ID,
				Name:      regPayload.Name,
				Age:       regPayload.Age,
				Bio:       "",
				Gender:    domain.Gender(regPayload.Gender),
				Location:  regPayload.Location,
				Interests: regPayload.Interests,
				Photos:    regPayload.Photos,
				IsPremium: false,
			}

			if regPayload.Tx != nil {
				if err := p.repo.CreateWithTx(ctx, regPayload.Tx, profile); err != nil {
					return fmt.Errorf("profile plugin failed to create profile with tx: %w", err)
				}
			} else {
				if err := p.repo.Create(ctx, profile); err != nil {
					return fmt.Errorf("profile plugin failed to create profile: %w", err)
				}
			}
			return nil
		})
	}

	return nil
}

func (p *ProfilePlugin) RegisterRoutes(apiGroup *echo.Group, authMiddleware echo.MiddlewareFunc) {
	userGroup := apiGroup.Group("/users", authMiddleware)
	p.handler.RegisterRoutes(userGroup)
}

func (p *ProfilePlugin) Repository() domain.ProfileRepository {
	return p.repo
}

func (p *ProfilePlugin) Service() domain.ProfileService {
	return p.service
}
