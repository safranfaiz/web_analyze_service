package main

import (
	"api/app"
	"api/configs"
	"fmt"
)

func main() {
    fmt.Println("Hello, World!")
    app.ApiRoutes(configs.GetAppConfig().ServerPort) 
}