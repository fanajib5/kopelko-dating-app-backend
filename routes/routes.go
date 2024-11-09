package routes

import (
	"kopelko-dating-app-backend/config"
	"kopelko-dating-app-backend/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(cfg *config.Config) *echo.Echo {
	// Initialize Echo
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = config.NewValidator()

	api := e.Group("/api")
	api.POST("/register", cfg.Controllers.Auth.RegisterUser)
	api.POST("/login", cfg.Controllers.Auth.LoginUser)

	user := api.Group("/users", utils.AuthMiddleware)
	profile := user.Group("/profiles", utils.AuthMiddleware)
	profile.GET("/me", cfg.Profile.ViewMyProfile)
	profile.GET("/random", cfg.Profile.RandomProfiles)
	user.POST("/swipe/:target_user_id", cfg.Swipe.SwipeHandler)
	user.POST("/subscribe", cfg.Subscribe.SubscribeHandler)

	return e
}
