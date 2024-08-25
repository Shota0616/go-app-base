package routes

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/Shota0616/go-sns/cmd/api/controllers"
	"github.com/Shota0616/go-sns/middleware" // ミドルウェアのパッケージ
	"github.com/gin-contrib/cors"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // ReactアプリケーションのURLを指定
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
		public.POST("/activate", controllers.Activate)
		public.POST("/login", controllers.Login)
		// サーバ側でトークンを管理するときは以下を追加
		// public.POST("/logout", controllers.Logout)
		public.GET("/ping", controllers.Ping)
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
