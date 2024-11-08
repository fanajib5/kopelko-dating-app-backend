package routes

import (
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/api")
	api.POST("/register", controllers.SignUp)
	// api.POST("/login", controllers.Login)

	// user := api.Group("/user")
	// user.GET("/profile", controllers.ViewProfile)
	// user.POST("/swipe", controllers.Swipe)
	// user.POST("/subscribe", controllers.Subscribe)
}
