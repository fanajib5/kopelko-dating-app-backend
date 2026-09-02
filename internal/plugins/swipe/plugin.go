package swipe

import (
	"kopelko-dating-app-backend/internal/core/plugin"
	profileRepo "kopelko-dating-app-backend/internal/plugins/profile/repository"
	swipeHttp "kopelko-dating-app-backend/internal/plugins/swipe/delivery/http"
	"kopelko-dating-app-backend/internal/plugins/swipe/domain"
	swipeRepo "kopelko-dating-app-backend/internal/plugins/swipe/repository"
	swipeUsecase "kopelko-dating-app-backend/internal/plugins/swipe/usecase"

	"github.com/labstack/echo/v4"
)

type SwipePlugin struct {
	repo    domain.SwipeRepository
	service domain.SwipeService
	handler *swipeHttp.SwipeHandler
}

func NewPlugin() plugin.Plugin {
	return &SwipePlugin{}
}

func (p *SwipePlugin) Name() string {
	return "swipe"
}

func (p *SwipePlugin) Version() string {
	return "1.0.0"
}

func (p *SwipePlugin) Init(appCtx *plugin.AppContext) error {
	p.repo = swipeRepo.NewSwipeRepository(appCtx.DBPool)
	profRepo := profileRepo.NewProfileRepository(appCtx.DBPool)

	p.service = swipeUsecase.NewSwipeUsecase(
		p.repo,
		profRepo,
		appCtx.Transactor,
		appCtx.Hooks,
		appCtx.Config.LimitSwipe,
	)
	p.handler = swipeHttp.NewSwipeHandler(p.service)

	return nil
}

func (p *SwipePlugin) RegisterRoutes(apiGroup *echo.Group, authMiddleware echo.MiddlewareFunc) {
	userGroup := apiGroup.Group("/users", authMiddleware)
	p.handler.RegisterRoutes(userGroup)
}

func (p *SwipePlugin) Service() domain.SwipeService {
	return p.service
}
