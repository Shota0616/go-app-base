package controllers

import (
	"net/http"
	// "time"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/Shota0616/go-sns/auth"
	"github.com/Shota0616/go-sns/models"
	"github.com/Shota0616/go-sns/config"
	"golang.org/x/crypto/bcrypt"
	"os"
	// "log"
)

// ユーザー登録を処理する関数
func Register(c *gin.Context) {
	url := os.Getenv("APP_URL") // アプリケーションのURLを環境変数から取得

	var input struct {
		Username string `json:"username"` // ユーザー名
		Password string `json:"password"` // パスワード
		Email    string `json:"email"`    // メールアドレス
	}

	// リクエストのJSONボディを構造体にバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
		return
	}

	// ユーザー情報をデータベースに保存
	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
		IsActive: false, // 初期状態では非アクティブ
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// JWTトークンとリフレッシュトークンを生成
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	// トークンとリフレッシュトークンをクライアントに返す
	c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refreshToken})

	// アクティベーションメールを送信
	subject := "仮登録のご案内"
	body := "以下のリンクをクリックして本登録を完了してください。\n" +
		url + "/auth/activate?token=" + token + "&email=" + user.Email

	if err := auth.SendEmail(user.Email, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}
}

// ユーザーのアクティベーションを処理する関数
func Activate(c *gin.Context) {
	var input struct {
		Email string `json:"email"` // メールアドレス
		Token string `json:"token"` // アクティベーション用トークン
	}

	// リクエストのJSONボディを構造体にバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Redisからトークンを取得
	val, err := config.RDB.Get(context.Background(), input.Email).Result()
	if err == redis.Nil || val != input.Token {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify token"})
		return
	}

	// ユーザー情報をデータベースから取得
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// ユーザーをアクティブにする
	user.IsActive = true
	config.DB.Save(&user)

	// Redisからトークンを削除
	config.RDB.Del(context.Background(), input.Email)

	c.JSON(http.StatusOK, gin.H{"message": "User activated. You can now log in."})
}

// ユーザーのログインを処理する関数
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`    // メールアドレス
		Password string `json:"password"` // パスワード
	}

	// リクエストのJSONボディを構造体にバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// メールアドレスでユーザーをデータベースから取得
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// パスワードの照合
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// アクティブなユーザーであることを確認
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is not activated"})
		return
	}

	// JWTトークンとリフレッシュトークンを生成
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	refreshtoken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	// トークンとリフレッシュトークンをクライアントに返す
	c.JSON(http.StatusOK, gin.H{"token": token, "refreshtoken": refreshtoken})
}

// トークンのリフレッシュを処理する関数
func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"` // リフレッシュトークン
	}

	// リクエストのJSONボディを構造体にバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// リフレッシュトークンの検証
	claims, err := auth.ValidateJWT(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// 新しいアクセストークンを生成
	newToken, err := auth.GenerateJWT(claims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new token"})
		return
	}

	// 新しいアクセストークンをクライアントに返す
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
