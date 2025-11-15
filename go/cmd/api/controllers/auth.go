package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"strings"
	"os"
	"encoding/json" // Added for JSON handling
	"log" // Added for logging
	"regexp" // Added for email validation

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-app-base/auth"
	"go-app-base/config"
	"go-app-base/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
	"github.com/mattn/go-sqlite3" // Added for SQLite error handling
	"errors" // Added for errors.As
)

// ランダムな文字列を生成するヘルパー関数
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// ユーザー登録を処理する関数
func Register(c *gin.Context) {
	var input struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// リクエストのJSONボディを構造体にバインド
    if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("ShouldBindJSON error: %v", err) // Debug log
		// 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "input_data_invalid"})})
        return
    }

	// メールアドレス形式バリデーション
	emailRegex := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	matched, _ := regexp.MatchString(emailRegex, input.Email)
	if !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_invalid"})})
		return
	}

	if len(input.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_min_length"})})
        return
    }

	// 既存ユーザーのバリデーションチェック
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		if !existingUser.IsActive {
			// 未バリデーションの場合は認証コード再発行・メール送信
			verificationCode, err := auth.GenerateVerificationCode()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_generate_verification_code"})})
				return
			}
			if err := config.RDB.Set(context.Background(), fmt.Sprintf("verification:%s", existingUser.Email), verificationCode, 10*time.Minute).Err(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "verification_code_save_failed"})})
				return
			}
			subject := "Your Verification Code"
			verifyURL := fmt.Sprintf("%s/auth/verify?email=%s&code=%s", os.Getenv("APP_URL"), existingUser.Email, verificationCode)
			body := fmt.Sprintf("Your verification code is: %s\n\nYou can verify your email here: %s", verificationCode, verifyURL)
			if err := auth.SendEmail(existingUser.Email, subject, body); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_send_failed"})})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_registration_success"})})
			return
		} else {
			// 既にバリデーション済みの場合はエラー
			errorMessage := config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_already_registered"})
			c.JSON(http.StatusConflict, gin.H{"error": errorMessage})
			return
		}
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_encryption_failed"})})
		return
	}

	var user models.User
	var result *gorm.DB
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		// ランダムなユーザー名を生成
		randomUsername := "user-" + generateRandomString(8)

		user = models.User{
			Username: randomUsername,
			Password: string(hashedPassword),
			Email:    input.Email,
			IsActive: false,
		}

		// ユーザーをデータベースに保存
		result = config.DB.Create(&user)
		if result.Error == nil {
			break // 成功したらループを抜ける
		}

		// エラーがユーザー名の重複によるものかチェック
		if strings.Contains(result.Error.Error(), "for key 'users.uni_users_username'") {
			// 重複している場合はリトライ
			continue
		} else {
			// その他のエラー
			break
		}
	}

	if result.Error != nil {
		var errorMessage string
		var sqliteErr sqlite3.Error
		if errors.As(result.Error, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			if strings.Contains(sqliteErr.Error(), "users.email") {
				errorMessage = config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_already_registered"})
			} else if strings.Contains(sqliteErr.Error(), "users.username") {
				errorMessage = config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "username_already_registered"})
			} else {
				errorMessage = config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_registration_failed"})
			}
		} else {
			errorMessage = config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_registration_failed"})
		}
		c.JSON(http.StatusConflict, gin.H{"error": errorMessage})
		return
	}


	// 認証コードを生成（6桁に統一）
	verificationCode, err := auth.GenerateVerificationCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_generate_verification_code"})})
		return
	}

	// Redisに認証コードを保存
	if err := config.RDB.Set(context.Background(), fmt.Sprintf("verification:%s", user.Email), verificationCode, 10*time.Minute).Err(); err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "verification_code_save_failed"})})
		return
	}

    // 認証コードをメールで送信
    subject := "Your Verification Code"
    verifyURL := fmt.Sprintf("%s/auth/verify?email=%s&code=%s", os.Getenv("APP_URL"), user.Email, verificationCode)
    body := fmt.Sprintf("Your verification code is: %s\n\nYou can verify your email here: %s", verificationCode, verifyURL)
    if err := auth.SendEmail(user.Email, subject, body); err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_send_failed"})})
        return
    }

	// 201 Created
	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_registration_success"})})
}

// ユーザーのログインを処理する関数
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`    // メールアドレス
		Password string `json:"password"` // パスワード
	}

	// リクエストのJSONボディを構造体にバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "input_data_invalid"})})
		return
	}

	// メールアドレスでユーザーをデータベースから取得
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
		return
	}

	// パスワードの照合
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_incorrect"})})
		return
	}

	// ユーザがアクティブかどうかを確認
	if (!user.IsActive) {
		// ユーザがアクティブでない場合、認証コードの入力画面にリダイレクトする。
		c.JSON(http.StatusSeeOther, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "account_not_activated_resend_verification"})})
        return
		// ここで処理を終了して、認証コードの入力画面にリダイレクトする
	}

	// JWTトークンとリフレッシュトークンを生成
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "could_not_generate_token"})})
		return
	}

	refreshtoken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "could_not_generate_refresh_token"})})
		return
	}

	// トークンとリフレシュトークンをクライアントに返す
	c.JSON(http.StatusOK, gin.H{"token": token, "refreshtoken": refreshtoken})
}


// トークンのリフレッシュを処理する関数
func RefreshToken(c *gin.Context) {
	    var input struct {
	        RefreshToken string `json:"refresh_token"` // リフレッシュトークン
	    }
	
	    // リクエストのJSONボディを構造体にバインド
	    if err := c.ShouldBindJSON(&input); err != nil {
	        c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "input_data_invalid"})})
	        return
	    }
	// リフレッシュトークンの検証
	claims, err := auth.ValidateJWT(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_refresh_token"})})
		return
	}

	// 新しいアクセストークンを生成
	newToken, err := auth.GenerateJWT(claims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "could_not_generate_new_token"})})
		return
	}

	// 新しいアクセストークンをクライアントに返す
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}


// EmailChangeData struct to unmarshal Redis JSON data
type EmailChangeData struct {
	UserID   uint   `json:"userID"`
	OldEmail string `json:"oldEmail"`
	Code     string `json:"code"`
}

// ユーザーの認証コードを検証する関数
func Verify(c *gin.Context) {
	var input struct {
		Email            string `json:"email"`
		VerificationCode string `json:"verificationCode"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "input_data_invalid"})})
		return
	}

	// Try to retrieve email change data from Redis
	redisKey := fmt.Sprintf("email_change_data:%s", input.Email)
	emailChangeDataJSON, err := config.RDB.Get(context.Background(), redisKey).Result()

	var emailChangeData EmailChangeData
	isEmailChangeFlow := false

	if err == nil {
		// Data found, it's an email change verification flow
		if err := json.Unmarshal([]byte(emailChangeDataJSON), &emailChangeData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email change data"})
			return
		}
		isEmailChangeFlow = true
	} else if err != redis.Nil {
		// Other Redis error
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "could_not_verify_code"})})
		return
	}

	if isEmailChangeFlow {
		// Email change verification flow
		if emailChangeData.Code != input.VerificationCode {
			c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_or_expired_verification_code"})})
			return
		}

		var user models.User
		if err := config.DB.First(&user, emailChangeData.UserID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
			return
		}

		user.Email = input.Email // Update to new email
		user.IsActive = true     // Ensure user is active
		if err := config.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_activate_user"})})
			return
		}

		// Delete email change data from Redis
		if err := config.RDB.Del(context.Background(), redisKey).Err(); err != nil {
			log.Printf("Failed to delete email change data from Redis for %s: %v", input.Email, err)
			// Log error but don't return, as the main task (user update) is done
		}

		c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_verified"})})
		return

	} else {
		// Regular registration verification flow
		val, err := config.RDB.Get(context.Background(), fmt.Sprintf("verification:%s", input.Email)).Result()
		if err == redis.Nil || val != input.VerificationCode {
			c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_or_expired_verification_code"})})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "could_not_verify_code"})})
			return
		}

		var user models.User
		if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
			return
		}

		user.IsActive = true
		if err := config.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_activate_user"})})
			return
		}

		// Delete verification code from Redis
		if err := config.RDB.Del(context.Background(), fmt.Sprintf("verification:%s", input.Email)).Err(); err != nil {
			log.Printf("Failed to delete verification code from Redis for %s: %v", input.Email, err)
		}

		// Delete resend count from Redis if exists
		resendKey := fmt.Sprintf("resend_count_%s", input.Email)
		if err := config.RDB.Del(context.Background(), resendKey).Err(); err != nil {
			log.Printf("Failed to delete resend count from Redis for %s: %v", input.Email, err)
		}

		c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_verified"})})
		return
	}
}

// 認証コードの再送を処理する関数
func ResendVerificationCode(c *gin.Context) {
    var input struct {
        Email string `json:"email"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "input_data_invalid"})})
        return
    }

	// ユーザーが存在するか確認
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
		return
	}

	// ユーザがアクティブかどうかを確認
	if (user.IsActive) {
		c.JSON(http.StatusConflict, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "account_already_activated"})})
		return
	}

    // 再送回数のチェック
    resendKey := fmt.Sprintf("resend_count_%s", input.Email)
    resendCount, err := config.RDB.Get(context.Background(), resendKey).Int()
    if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_check_resend_count"})})
        return
    }

	// 再送回数が3回を超えた場合、エラーメッセージを返す
    if resendCount >= 3 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resend_limit_reached"})})
        return
    }

    // 認証コードを生成（6桁に統一）
    verificationCode, err := auth.GenerateVerificationCode()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_generate_verification_code"})})
        return
    }

	// redisに保存されている認証コードを削除
	if err := config.RDB.Del(context.Background(), input.Email).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "verification_code_delete_failed"})})
		return
	}

    // Redisに認証コードを保存
    if err := config.RDB.Set(context.Background(), input.Email, verificationCode, 10*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "verification_code_save_failed"})})
        return
    }

    // 認証コードをメールで送信
    subject := "Your Verification Code"
    body := fmt.Sprintf("Your verification code is: %s", verificationCode)
    if err := auth.SendEmail(input.Email, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_send_failed"})})
        return
    }

    // 再送回数をインクリメント
    if err := config.RDB.Incr(context.Background(), resendKey).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_increment_resend_count"})})
        return
    }
    // 再送回数の有効期限を12時間に設定
    if err := config.RDB.Expire(context.Background(), resendKey, 12*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_set_resend_count_expiration"})})
        return
    }

	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "verification_code_resent"})})
}


// パスワード再設定リクエストを処理する関数
func RequestPasswordReset(c *gin.Context) {
    var input struct {
        Email string `json:"email"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "input_data_invalid"})})
        return
    }

    var user models.User
    if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
        return
    }

	// ユーザがアクティブかどうかを確認
	if (!user.IsActive) {
		// ユーザがアクティブでない場合、認証コードの入力画面にリダイレクトする。
		c.JSON(http.StatusSeeOther, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "account_not_activated_resend_verification"})})
		return
		// ここで処理を終了して、認証コードの入力画面にリダイレクトする
	}

	// 再送回数のチェック
	resendKey := fmt.Sprintf("resend_count_%s", user.Email)
	resendCount, err := config.RDB.Get(context.Background(), resendKey).Int()
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_check_resend_count"})})
		return
	}

	// 再送回数が3回を超えた場合、エラーメッセージを返す
    if resendCount >= 3 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resend_limit_reached"})})
        return
    }


    // トークンを生成,useridをキーにしてRedisに保存
	token, err := auth.GenerateJWT(user.ID)
    if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "token_generation_failed"})})
        return
    }

    // Redisにトークンを保存
    if err := config.RDB.Set(context.Background(), user.Email, token, 10*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "token_save_failed"})})
        return
    }

    // トークンをメールで送信
    resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", os.Getenv("APP_URL"), token)
	subject := "Password Reset Link"
	body := fmt.Sprintf("The link to reset your password is: %s", resetURL)
    if err := auth.SendEmail(user.Email, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_send_failed"})})
        return
    }

    // 再送回数をインクリメント
    if err := config.RDB.Incr(context.Background(), resendKey).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_increment_resend_count"})})
        return
    }
    // 再送回数の有効期限を12時間に設定
    if err := config.RDB.Expire(context.Background(), resendKey, 12*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_set_resend_count_expiration"})})
        return
    }


	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset_link_sent"})})
}


// パスワード再設定を処理する関数
func ResetPassword(c *gin.Context) {
    var input struct {
        Token       string `json:"token"`
        NewPassword string `json:"newPassword"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "input_data_invalid"})})
        return
    }

    // トークンを検証し、クレームを取得
    claims, err := auth.ValidateJWT(input.Token)
    if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_or_expired_token"})})
        return
    }

	// クレームからユーザーIDを取得
	userID := claims.ID

    // トークンのクレームから取得したユーザーIDでユーザーを取得
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_not_found"})})
		return
	}

    // 新しいパスワードをハッシュ化
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
    if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_encryption_failed"})})
        return
    }

    // パスワードを更新
    user.Password = string(hashedPassword)
    if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_update_failed"})})
        return
    }

    // Redisからトークンを削除
    if err := config.RDB.Del(context.Background(), input.Token).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "token_delete_failed"})})
        return
    }

	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset_success"})})
}