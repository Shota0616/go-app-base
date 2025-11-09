package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"go-app-base/config"
	"go-app-base/models"
	"go-app-base/auth"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetUser(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	log.Println(tokenStr)

	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "authorization_token_not_provided"})})
		return
	}

	claims, err := auth.ValidateJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_token"})})
		return
	}

	userID := claims.ID

	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"active":   user.IsActive,
	})
}

func UpdateUsername(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
	}

	// リクエストのJSONをバインドし、バリデーション
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// middlewareからユーザーIDを取得
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		return
	}

	// ユーザーを検索
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// ユーザー名を更新
	user.Username = input.Username
	if err := config.DB.Save(&user).Error; err != nil {
		// 重複エラーのハンドリング
		if strings.Contains(err.Error(), "for key 'users.uni_users_username'") {
			c.JSON(http.StatusConflict, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "username_already_registered"})})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}


func UpdateUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "authorization_token_not_provided"})})
		return
	}

	claims, err := auth.ValidateJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_token"})})
		return
	}

	userID := claims.ID

	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
		return
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_encryption_failed"})})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_update_failed"})})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_updated_successfully"})})
}

func DeleteUser(c *gin.Context) {
    tokenStr := c.GetHeader("Authorization")
    if tokenStr == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "authorization_token_not_provided"})})
        return
    }

    claims, err := auth.ValidateJWT(tokenStr)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_token"})})
        return
    }

    userID := claims.ID

    var user models.User
    if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
        return
    }

    if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_deletion_failed"})})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_deleted_successfully"})})
}
