package app

import (
	"api/app/handler"
	"api/configs"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ApiRoutes(port string) {
	routes := gin.Default()

	api := CorsConfig(routes)
	{
		api.GET("/analyze", handler.WebPageExecutorHandler)
	}

	err := routes.Run(port)
	if err != nil {
		log.Fatal("Failed to start the server", err)
	}
}

// CORS config
func CorsConfig(routes *gin.Engine) *gin.RouterGroup {
	routes.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := routes.Group(configs.GetConfig().BasePath + configs.GetConfig().ApiVersion)
	return api
}
