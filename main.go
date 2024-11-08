package main

import (
	"kopelko-dating-app-backend/routes"
	"kopelko-dating-app-backend/utils"
)

func main() {
	e := routes.SetupRoutes()

	utils.NewLogger()
	// e.Logger.Fatal(e.Start(":8080"))
	utils.Logger.LogInfo().Msg(e.Start(":8080").Error())
}
