package main

import (
	"api/app"
	"api/configs"
	"log"
)

func main() {

	// Load application configuration
	cfg := configs.GetConfig()
	log.Println("Starting API server on port", cfg.ServerPort)

	// Start the server
	app.ApiRoutes(cfg.ServerPort)
}
