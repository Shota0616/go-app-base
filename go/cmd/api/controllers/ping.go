package controllers

import (
	"github.com/gin-gonic/gin"
	// "github.com/go-redis/redis/v8"
	// "go-app-base/config"
	// "go-app-base/models"
	// "go-app-base/auth"
	// "golang.org/x/crypto/bcrypt"
)

func Ping(c *gin.Context) {
	// JSONレスポンスを返す
	c.JSON(200, gin.H{
		"message": "Hello World! Pong",
	})
}