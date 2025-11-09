package routes

import (
	"time"
	"github.com/gin-gonic/gin"
	"go-app-base/cmd/api/controllers"
	"go-app-base/middleware" // ミドルウェアのパッケージ
	"github.com/gin-contrib/cors"
	"os"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173", "http://127.0.0.1:3000", os.Getenv("APP_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// パブリックルート
	public := router.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/verify", controllers.Verify)
		public.POST("/login", controllers.Login)
		public.POST("/request-password-reset", controllers.RequestPasswordReset) // パスワード再設定リクエストのエンドポイントを追加
		public.POST("/resend-verification-code", controllers.ResendVerificationCode) // メール認証コード再送のエンドポイントを追加
		public.POST("/reset-password", controllers.ResetPassword) // パスワード再設定のエンドポイントを追加
		public.GET("/ping", controllers.Ping)
		public.GET("/db-check", controllers.DBCheck)
		// サーバ側でトークンを管理するときは以下を追加
		// public.POST("/logout", controllers.Logout)
	}

	// 認証が必要なルート
	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		// protected.GET("/mypage", controllers.GetMyPage) // マイページ
		protected.GET("/getuser", controllers.GetUser) // ユーザー情報取得
		protected.PUT("/user/username", controllers.UpdateUsername) // ユーザー名更新
		// その他の保護されたルート
	}

	return router
}
