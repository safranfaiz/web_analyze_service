package app

import (
	"api/app/handler"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(port string) {
	routes := gin.Default()
	routes.GET("/test", handler.TestHandler)
	
	routes.Run(port)
}