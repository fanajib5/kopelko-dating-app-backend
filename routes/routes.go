package routes

import (
	"kopelko-dating-app-backend/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes() *echo.Echo {
	// Initialize config
	cfg := config.New()

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = config.NewValidator()

	api := e.Group("/api")
	api.POST("/register", cfg.Controllers.Auth.RegisterUser)
	api.POST("/login", cfg.Controllers.Auth.LoginUser)

	// user := api.Group("/user")
	// user.GET("/profile", controllers.ViewProfile)
	// user.POST("/swipe", controllers.Swipe)
	// user.POST("/subscribe", controllers.Subscribe)

	return e
}
