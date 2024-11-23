package routes

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/Shota0616/go-sns/cmd/api/controllers"
	"github.com/Shota0616/go-sns/middleware" // ミドルウェアのパッケージ
	"github.com/gin-contrib/cors"
	"os"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("APP_URL"),"http://localhost:5173"}, // .envファイルに設定したAPP_URLを使用
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
		// サーバ側でトークンを管理するときは以下を追加
		// public.POST("/logout", controllers.Logout)
	}

	// 認証が必要なルート
	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		// protected.GET("/mypage", controllers.GetMyPage) // マイページ
		protected.GET("/getuser", controllers.GetUser) // ユーザー情報取得
		// その他の保護されたルート
	}

	return router
}
