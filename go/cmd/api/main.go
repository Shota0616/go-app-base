package main

import (
	// "github.com/gin-gonic/gin"
	"github.com/Shota0616/go-sns/config"
	"github.com/Shota0616/go-sns/routes"
	// "log"
)

func main() {
	config.ConnectDatabase()
	config.MigrateDatabase()
	config.ConnectRedis()

	router := routes.SetupRouter()
	router.Run(":8080")
}