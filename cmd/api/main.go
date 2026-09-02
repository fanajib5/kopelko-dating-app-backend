package main

import (
	"log"

	identityHttp "kopelko-dating-app-backend/internal/modules/identity/delivery/http"
	identityRepo "kopelko-dating-app-backend/internal/modules/identity/repository"
	identityUsecase "kopelko-dating-app-backend/internal/modules/identity/usecase"

	profileHttp "kopelko-dating-app-backend/internal/modules/profile/delivery/http"
	profileRepo "kopelko-dating-app-backend/internal/modules/profile/repository"
	profileUsecase "kopelko-dating-app-backend/internal/modules/profile/usecase"

	subscriptionHttp "kopelko-dating-app-backend/internal/modules/subscription/delivery/http"
	subscriptionRepo "kopelko-dating-app-backend/internal/modules/subscription/repository"
	subscriptionUsecase "kopelko-dating-app-backend/internal/modules/subscription/usecase"

	swipeHttp "kopelko-dating-app-backend/internal/modules/swipe/delivery/http"
	swipeRepo "kopelko-dating-app-backend/internal/modules/swipe/repository"
	swipeUsecase "kopelko-dating-app-backend/internal/modules/swipe/usecase"

	"kopelko-dating-app-backend/internal/platform/config"
	"kopelko-dating-app-backend/internal/platform/database"
	"kopelko-dating-app-backend/internal/platform/middleware"
	"kopelko-dating-app-backend/internal/platform/token"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Println("Starting dating app backend modular monolith...")

	// 1. Config
	cfg := config.Load()

	// 2. Database Connection & Transactor
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer dbPool.Close()

	transactor := database.NewTransactor(dbPool)

	// 3. Platform Services & Middleware
	tokenSvc := token.NewJWTService(cfg.JWTSecret)
	customValidator := middleware.NewValidator()
	authMiddleware := middleware.AuthMiddleware(tokenSvc)

	// 4. Repositories
	userRepository := identityRepo.NewUserRepository(dbPool)
	profileRepository := profileRepo.NewProfileRepository(dbPool)
	subscriptionRepository := subscriptionRepo.NewSubscriptionRepository(dbPool)
	swipeRepository := swipeRepo.NewSwipeRepository(dbPool)

	// 5. Usecases / Domain Services
	subscriptionService := subscriptionUsecase.NewSubscriptionUsecase(subscriptionRepository)
	profileService := profileUsecase.NewProfileUsecase(profileRepository, subscriptionService, cfg.LimitSwipe)
	swipeService := swipeUsecase.NewSwipeUsecase(swipeRepository, subscriptionService, profileRepository, transactor, cfg.LimitSwipe)
	identityService := identityUsecase.NewIdentityUsecase(userRepository, profileRepository, tokenSvc, transactor)

	// 6. HTTP Handlers
	identityHandler := identityHttp.NewIdentityHandler(identityService)
	profileHandler := profileHttp.NewProfileHandler(profileService)
	subscriptionHandler := subscriptionHttp.NewSubscriptionHandler(subscriptionService)
	swipeHandler := swipeHttp.NewSwipeHandler(swipeService)

	// 7. Echo Router Setup
	e := echo.New()
	e.HideBanner = true
	e.Validator = customValidator
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	api := e.Group("/api")
	identityHandler.RegisterRoutes(api)

	// Protected Routes
	userGroup := api.Group("/users", authMiddleware)
	profileHandler.RegisterRoutes(userGroup)
	subscriptionHandler.RegisterRoutes(userGroup)
	swipeHandler.RegisterRoutes(userGroup)

	// 8. Graceful Server Startup & Shutdown
	go func() {
		addr := ":" + cfg.APIPort
		log.Printf("Server listening on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
