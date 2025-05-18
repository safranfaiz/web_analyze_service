package main

import (
	"api/app"
	"api/configs"
)

func main() {

	// load the configuration for application
	configs.GetConfig().LoadConfig()

	// start the server
	app.ApiRoutes(configs.GetConfig().ServerPort)
}
