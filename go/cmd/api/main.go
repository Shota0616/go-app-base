package main

import (
	// "github.com/gin-gonic/gin"
	"go-app-base/config"
	"go-app-base/routes"
	// "log"
)

func main() {
	config.InitI18n()
	config.ConnectDatabase()
	config.MigrateDatabase()
	config.ConnectRedis()

	router := routes.SetupRouter()
	router.Run(":8080")
}