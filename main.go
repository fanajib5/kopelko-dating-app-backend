package main

import (
	"kopelko-dating-app-backend/utils"
	"log"
)

func main() {
	// Initialize database connection
	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("could not set up database: %v", err)
	}
}
