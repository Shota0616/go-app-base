package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-app-base/config"
)

// DBCheck はデータベースへの接続を確認します。
func DBCheck(c *gin.Context) {
	sqlDB, err := config.DB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get database instance",
		})
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to connect to database",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Database connection is healthy",
	})
}
