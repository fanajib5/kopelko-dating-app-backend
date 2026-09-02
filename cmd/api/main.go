package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	_ "kopelko-dating-app-backend/docs"
	"kopelko-dating-app-backend/internal/platform/config"
	"kopelko-dating-app-backend/internal/platform/database"
	pkgLogger "kopelko-dating-app-backend/internal/platform/logger"
	"kopelko-dating-app-backend/internal/platform/middleware"
	"kopelko-dating-app-backend/internal/platform/migrator"
	"kopelko-dating-app-backend/internal/platform/token"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Kopelko Dating App API
// @version 1.0
// @description Modular Monolith Dating App Backend REST API built with Echo & pgx
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	seedFlag := flag.Bool("seed", false, "Run database seeders")
	flag.Parse()

	// 1. Config & Structured Logger
	cfg := config.Load()
	pkgLogger.Init(cfg.AppEnv)
	slog.Info("Starting dating app backend modular monolith...", "env", cfg.AppEnv, "port", cfg.APIPort)

	// 2. Database Connection & Transactor
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Database connection pool failed", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Run CLI Migrations / Seeders if flagged
	if *migrateFlag {
		if err := migrator.RunMigration(context.Background(), dbPool, "databases/migrations/schema.sql"); err != nil {
			slog.Error("Migration failed", "error", err)
			os.Exit(1)
		}
		if !*seedFlag {
			slog.Info("Migration completed successfully")
			return
		}
	}

	if *seedFlag {
		if err := migrator.RunSeeder(context.Background(), dbPool, "databases/seeders/seeder.sql"); err != nil {
			slog.Error("Seeding failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Seeding completed successfully")
		return
	}

	transactor := database.NewTransactor(dbPool)

	// 3. Platform Services & Middleware
	tokenSvc := token.NewJWTService(cfg.JWTSecret)
	customValidator := middleware.NewValidator()
	authMiddleware := middleware.AuthMiddleware(tokenSvc)
	securityHeadersMiddleware := middleware.SecurityHeadersMiddleware()
	authRateLimiter := middleware.NewAuthRateLimiter(10, 1*time.Minute)

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
	e.Use(echoMiddleware.Recover())
	e.Use(securityHeadersMiddleware)
	e.Use(middleware.RequestIDAndLoggingMiddleware())

	// Swagger documentation UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	api := e.Group("/api")

	// Public Auth Routes protected with rate limiter
	authGroup := api.Group("", authRateLimiter)
	identityHandler.RegisterRoutes(authGroup)

	// Protected Routes
	userGroup := api.Group("/users", authMiddleware)
	profileHandler.RegisterRoutes(userGroup)
	subscriptionHandler.RegisterRoutes(userGroup)
	swipeHandler.RegisterRoutes(userGroup)

	// 8. Graceful Server Startup & Shutdown
	go func() {
		addr := fmt.Sprintf(":%s", cfg.APIPort)
		slog.Info("HTTP Server listening", "address", addr, "swagger", fmt.Sprintf("http://localhost:%s/swagger/index.html", cfg.APIPort))
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			slog.Error("Server encountered fatal error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited cleanly")
}
