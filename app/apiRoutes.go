package app

import (
	"api/app/handler"
	"api/configs"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ApiRoutes sets up and starts the server on the given port.
func ApiRoutes(port string) {
	router := gin.Default()

	// Apply CORS and define API route group
	apiGroup := applyCORS(router)

	// Register route handlers
	apiGroup.GET("/analyze", handler.WebPageExecutorHandler)

	// Start the server
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}

// applyCORS sets up CORS middleware and returns the base API group.
func applyCORS(router *gin.Engine) *gin.RouterGroup {
	cfg := configs.GetConfig()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	basePath := cfg.BasePath + cfg.ApiVersion
	return router.Group(basePath)
}
