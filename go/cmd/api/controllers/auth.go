package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"strings"
	"os"


	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/Shota0616/go-sns/auth"
	"github.com/Shota0616/go-sns/config"
	"github.com/Shota0616/go-sns/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// ユーザー登録を処理する関数
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// リクエストのJSONボディを構造体にバインド
    if err := c.ShouldBindJSON(&input); err != nil {
		// 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_input"})})
        return
    }

	// パスワードをハッシュ化, hashedPasswordにはハッシュ化されたパスワードが格納される, errにはエラーが格納される
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	// ハッシュ化に失敗した場合jsonでエラーメッセージを返す
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_encryption_failed"})})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
		IsActive: false,
	}

	// ユーザーをデータベースに保存, エラーが発生したらエラーメッセージをjsonで返す
	if err := config.DB.Create(&user).Error; err != nil {
		var errorMessage string
		switch {
		case strings.Contains(err.Error(), "for key 'users.uni_users_email'"):
			errorMessage = config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_already_registered"})
		case strings.Contains(err.Error(), "for key 'users.uni_users_username'"):
			errorMessage = config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "username_already_registered"})
		default:
			errorMessage = config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_registration_failed"})
		}
		// 409 Conflict
		c.JSON(http.StatusConflict, gin.H{"error": errorMessage})
		// c.JSON(http.StatusConflict, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: errorMessage})})
		return
	}


	// 4桁の認証コードを生成
	rand.Seed(time.Now().UnixNano())
	verificationCode := fmt.Sprintf("%04d", rand.Intn(10000))

	// Redisに認証コードを保存
	if err := config.RDB.Set(context.Background(), user.Email, verificationCode, 10*time.Minute).Err(); err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "verification_code_save_failed"})})
		return
	}

    // 認証コードをメールで送信
    subject := "Your Verification Code"
    body := fmt.Sprintf("Your verification code is: %s", verificationCode)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_input"})})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "incorrect_password"})})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_input"})})
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


// ユーザーの認証コードを検証する関数
func Verify(c *gin.Context) {
	var input struct {
		Email            string `json:"email"`
		VerificationCode string `json:"verificationCode"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_input"})})
		return
	}

	// Redisから認証コードを取得
	val, err := config.RDB.Get(context.Background(), input.Email).Result()
	if err == redis.Nil || val != input.VerificationCode {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_or_expired_verification_code"})})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "could_not_verify_code"})})
		return
	}

	// ユーザをアクティブにする
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

	// Redisから認証コードを削除
	if err := config.RDB.Del(context.Background(), input.Email).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "verification_code_delete_failed"})})
		return
	}

	// もし再送回数のキーが存在していたら削除
	resendKey := fmt.Sprintf("resend_count_%s", input.Email)
	if err := config.RDB.Del(context.Background(), resendKey).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "failed_to_delete_resend_count"})})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_verified",})})
}

// 認証コードの再送を処理する関数
func ResendVerificationCode(c *gin.Context) {
    var input struct {
        Email string `json:"email"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_input"})})
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

    // 4桁の認証コードを生成
    rand.Seed(time.Now().UnixNano())
    verificationCode := fmt.Sprintf("%04d", rand.Intn(10000))

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
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_input"})})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_input"})})
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