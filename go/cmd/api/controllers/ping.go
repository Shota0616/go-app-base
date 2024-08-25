package controllers

import (
	"github.com/gin-gonic/gin"
	// "github.com/go-redis/redis/v8"
	// "github.com/Shota0616/go-sns/config"
	// "github.com/Shota0616/go-sns/models"
	// "github.com/Shota0616/go-sns/auth"
	// "golang.org/x/crypto/bcrypt"
)

func Ping(c *gin.Context) {
	// JSONレスポンスを返す
	c.JSON(200, gin.H{
		"message": "Hello World! Pong",
	})
}