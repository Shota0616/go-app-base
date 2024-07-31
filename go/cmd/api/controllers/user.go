package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Shota0616/go-sns/config"
	"github.com/Shota0616/go-sns/models"
	"github.com/Shota0616/go-sns/auth"
	"golang.org/x/crypto/bcrypt"
	// "context"
)

func GetUser(c *gin.Context) {
	// username := c.Param("username")

	var user models.User
	// if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"active":   user.IsActive,
	})
}

func UpdateUser(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// JWTトークンからユーザー情報を取得
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
		return
	}

	claims, err := auth.ValidateJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", claims.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
			return
		}
		user.Password = string(hashedPassword)
	}

	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	// JWTトークンからユーザー情報を取得
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
		return
	}

	claims, err := auth.ValidateJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", claims.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
