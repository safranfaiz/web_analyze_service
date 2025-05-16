package app

import (
	"api/app/handler"
	"api/configs"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(port string) {
	routes := gin.Default()	
	api := routes.Group(configs.GetConfig().BasePath + configs.GetConfig().ApiVersion) 
	{
		api.GET("/analyze", handler.WebPageExecutorHandler)
	}	
	routes.Run(port)
}