package main

import (
	"kopelko-dating-app-backend/config"
	"kopelko-dating-app-backend/routes"
)

func main() {
	// Initialize config
	cfg := config.New()

	e := routes.SetupRoutes(cfg)
	e.Logger.Fatal(e.Start(":" + cfg.APIPort))
}
