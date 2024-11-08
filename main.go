package main

import "kopelko-dating-app-backend/routes"

func main() {
	e := routes.SetupRoutes()
	e.Logger.Fatal(e.Start(":8080"))
}
