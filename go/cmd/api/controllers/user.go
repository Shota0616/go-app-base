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
	"time" // Added for email sending
	"fmt" // Added for fmt.Sprintf
	"encoding/json" // Added for JSON handling
	"errors" // Added for errors.As
	"github.com/mattn/go-sqlite3" // Added for SQLite error handling
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

	log.Printf("Attempting to update username for userID: %v with new username: %s", userID, input.Username)

	// ユーザーを検索
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		log.Printf("Error finding user %v: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// ユーザー名を更新
	user.Username = input.Username
	if err := config.DB.Save(&user).Error; err != nil {
		log.Printf("Error saving user %v with new username %s: %v", userID, input.Username, err)
		
		// Check for duplicate username error (MySQL or SQLite)
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusConflict, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "username_already_registered"})})
			return
		}
		
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			if strings.Contains(sqliteErr.Error(), "users.username") {
				c.JSON(http.StatusConflict, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "username_already_registered"})})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_update_username"})})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "username_updated_successfully"})}) // Added translation key
}

func UpdateEmail(c *gin.Context) {
	var input struct {
		NewEmail string `json:"newEmail" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
		return
	}

	// If new email is the same as current email, just return success
	if user.Email == input.NewEmail {
		c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_already_current"})})
		return
	}

	// Check if the new email is already registered by another user
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.NewEmail).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_already_registered"})})
		return
	}

	// Generate verification code
	verificationCode, err := auth.GenerateVerificationCode()
	if err != nil {
		log.Printf("Failed to generate verification code for email change for user %s: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_generate_verification_code"})})
		return
	}

	// Store new email and user ID temporarily in Redis with the verification code
	// Key: "email_change_data:{newEmail}"
	// Value: JSON string {"userID": userID, "oldEmail": "...", "code": "..."}
	emailChangeData := map[string]interface{}{
		"userID":   userID,
		"oldEmail": user.Email,
		"code":     verificationCode,
	}
	emailChangeDataJSON, err := json.Marshal(emailChangeData)
	if err != nil {
		log.Printf("Failed to marshal email change data for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process email change"})
		return
	}

	redisKey := fmt.Sprintf("email_change_data:%s", input.NewEmail)
	if err := config.RDB.Set(c, redisKey, emailChangeDataJSON, time.Minute*10).Err(); err != nil { // 10 minutes expiration
		log.Printf("Failed to store email change data for user %d in Redis: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store email change data"})
		return
	}

	// Send verification email to the NEW email address
	if err := auth.SendVerificationEmail(input.NewEmail, verificationCode); err != nil {
		log.Printf("Failed to send verification email for email change to %s: %v", input.NewEmail, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_send_email"})})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_updated_successfully_please_verify"})})
}

func UpdatePassword(c *gin.Context) {
	var input struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "incorrect_current_password"})})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_encryption_failed"})})
		return
	}
	user.Password = string(hashedPassword)

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_update_failed"})})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_updated_successfully"})})
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
