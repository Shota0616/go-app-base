package routes

import (
	"time"
	"github.com/gin-gonic/gin"
	"go-app-base/cmd/api/controllers"
	"go-app-base/config"
	"go-app-base/middleware"
	"github.com/gin-contrib/cors"
)

func SetupRouter() *gin.Engine {
	// 本番環境ではリリースモードに設定
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.Default()

	// CORS設定（環境変数から取得）
	allowedOrigins := config.GetAllowedOrigins()
	if appURL := config.GetEnv("APP_URL", ""); appURL != "" {
		allowedOrigins = append(allowedOrigins, appURL)
	}
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
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
		public.POST("/request-password-reset", controllers.RequestPasswordReset)
		public.POST("/resend-verification-code", controllers.ResendVerificationCode)
		public.POST("/reset-password", controllers.ResetPassword)
		public.GET("/ping", controllers.Ping)
		public.GET("/db-check", controllers.DBCheck)
	}

	// 認証が必要なルート
	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/getuser", controllers.GetUser)
		protected.PUT("/user/username", controllers.UpdateUsername)
		protected.PUT("/user/email", controllers.UpdateEmail)
		protected.PUT("/user/password", controllers.UpdatePassword)
	}

	return router
}
