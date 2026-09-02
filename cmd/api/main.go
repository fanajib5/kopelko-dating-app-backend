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

	_ "kopelko-dating-app-backend/docs"
	"kopelko-dating-app-backend/internal/core/hook"
	identityHttp "kopelko-dating-app-backend/internal/core/identity/delivery/http"
	identityRepo "kopelko-dating-app-backend/internal/core/identity/repository"
	identityUsecase "kopelko-dating-app-backend/internal/core/identity/usecase"
	"kopelko-dating-app-backend/internal/core/plugin"
	"kopelko-dating-app-backend/internal/platform/config"
	"kopelko-dating-app-backend/internal/platform/database"
	pkgLogger "kopelko-dating-app-backend/internal/platform/logger"
	"kopelko-dating-app-backend/internal/platform/middleware"
	"kopelko-dating-app-backend/internal/platform/migrator"
	"kopelko-dating-app-backend/internal/platform/token"
	profilePlugin "kopelko-dating-app-backend/internal/plugins/profile"
	subscriptionPlugin "kopelko-dating-app-backend/internal/plugins/subscription"
	swipePlugin "kopelko-dating-app-backend/internal/plugins/swipe"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Kopelko Dating App API
// @version 1.0
// @description Modular Extensible Dating App Backend Platform (Core + WordPress-like Plugins & Hooks)
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
	slog.Info("Starting Kopelko Platform Core...", "env", cfg.AppEnv, "port", cfg.APIPort)

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

	// 3. Platform Services, Middleware & Hook Manager
	hookManager := hook.NewHookManager()
	tokenSvc := token.NewJWTService(cfg.JWTSecret)
	customValidator := middleware.NewValidator()
	authMiddleware := middleware.AuthMiddleware(tokenSvc)
	securityHeadersMiddleware := middleware.SecurityHeadersMiddleware()
	authRateLimiter := middleware.NewAuthRateLimiter(10, 1*time.Minute)

	// 4. Core Modules (Identity & Authentication)
	userRepository := identityRepo.NewUserRepository(dbPool)
	identityService := identityUsecase.NewIdentityUsecase(userRepository, tokenSvc, transactor, hookManager)
	identityHandler := identityHttp.NewIdentityHandler(identityService)

	// 5. Plugin Architecture Setup
	appCtx := &plugin.AppContext{
		Context:    context.Background(),
		DBPool:     dbPool,
		Transactor: transactor,
		Hooks:      hookManager,
		Config:     cfg,
		Logger:     slog.Default(),
	}

	pluginMgr := plugin.NewPluginManager(appCtx)

	// Register Plugins (subscription, profile, swipe)
	pluginMgr.Register(subscriptionPlugin.NewPlugin())
	pluginMgr.Register(profilePlugin.NewPlugin())
	pluginMgr.Register(swipePlugin.NewPlugin())

	// Initialize all plugins & register hooks
	if err := pluginMgr.InitAll(); err != nil {
		slog.Error("Failed to initialize plugins", "error", err)
		os.Exit(1)
	}

	// 6. Echo Router Setup
	e := echo.New()
	e.HideBanner = true
	e.Validator = customValidator
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.Recover())
	e.Use(securityHeadersMiddleware)
	e.Use(middleware.RequestIDAndLoggingMiddleware())

	// Swagger documentation UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Healthcheck Probes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "alive",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	e.GET("/health/ready", func(c echo.Context) error {
		pingCtx, pingCancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer pingCancel()

		if err := dbPool.Ping(pingCtx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]any{
				"status":   "unready",
				"database": "disconnected",
				"error":    err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"status":   "ready",
			"database": "connected",
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	})

	api := e.Group("/api")

	// Core Auth Routes (protected with rate limiter)
	authGroup := api.Group("", authRateLimiter)
	identityHandler.RegisterRoutes(authGroup)

	// Mount Plugin Routes
	pluginMgr.RegisterAllRoutes(api, authMiddleware)

	// 7. Graceful Server Startup & Shutdown
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
