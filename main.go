package main

import (
	"log"

	"kopelko-dating-app-backend/config"
	"kopelko-dating-app-backend/routes"
)

func main() {
	log.Println("Starting the application...")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize config
	cfg := config.New()

	e := routes.SetupRoutes(cfg)
	e.Logger.Fatal(e.Start(":" + cfg.APIPort))
}
