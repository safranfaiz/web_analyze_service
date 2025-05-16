package app

import (
	"api/app/handler"
	"api/configs"
	"log"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(port string) {
	routes := gin.Default()	
	api := routes.Group(configs.GetConfig().BasePath + configs.GetConfig().ApiVersion) 
	{
		api.GET("/analyze", handler.WebPageExecutorHandler)
	}
		
	err := routes.Run(port)
	if err != nil {
		log.Fatal("Failed to start the server", err)
	}
}