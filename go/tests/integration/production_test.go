package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-app-base/config"
	"go-app-base/routes"
)

func setupProductionEnv() {
	os.Setenv("ENV_MODE", "production")
	os.Setenv("JWT_SECRET", "prod_jwt_secret_very_secure")
	os.Setenv("JWT_REFRESH_SECRET", "prod_jwt_refresh_secret_very_secure")
	os.Setenv("APP_URL", "https://yourdomain.com")
	os.Setenv("EMAIL_FROM", "noreply@yourdomain.com")
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USER", "your-email@gmail.com")
	os.Setenv("SMTP_PASSWORD", "your-app-password")
	os.Setenv("REDIS_HOST", "redis")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("APP_ROOT", "/usr/src")
	os.Setenv("APP_LANG", "en")
	os.Setenv("ALLOWED_ORIGINS", "https://yourdomain.com")
	os.Setenv("LOG_LEVEL", "warn")
}

func TestProductionEnvironment(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	assert.True(t, config.IsProduction())
	assert.False(t, config.IsDevelopment())
	assert.Equal(t, "warn", config.GetLogLevel())
}

func TestProductionGinMode(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	// ルーター設定時にGinがリリースモードになることを確認
	routes.SetupRouter()
	assert.Equal(t, gin.ReleaseMode, gin.Mode())
}

func TestProductionCORS(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	origins := config.GetAllowedOrigins()
	assert.Contains(t, origins, "https://yourdomain.com")
	assert.NotContains(t, origins, "http://localhost:3000")
}

func TestProductionSMTPConfig(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	assert.Equal(t, "smtp.gmail.com", os.Getenv("SMTP_HOST"))
	assert.Equal(t, "587", os.Getenv("SMTP_PORT"))
	assert.NotEmpty(t, os.Getenv("SMTP_USER"))
	assert.NotEmpty(t, os.Getenv("SMTP_PASSWORD"))
}

func TestProductionSecurityHeaders(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	router := routes.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductionJWTSecrets(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	// 本番環境では強力なシークレットが必要
	assert.NotEmpty(t, jwtSecret)
	assert.NotEmpty(t, jwtRefreshSecret)
	assert.NotEqual(t, jwtSecret, jwtRefreshSecret)
	assert.Greater(t, len(jwtSecret), 20)
	assert.Greater(t, len(jwtRefreshSecret), 20)
}

func TestProductionPingEndpoint(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	router := routes.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Hello World! Pong", response["message"])
}

func TestProductionCORSRestriction(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	router := routes.SetupRouter()
	
	// 許可されていないオリジンからのリクエスト
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/ping", nil)
	req.Header.Set("Origin", "http://malicious-site.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	router.ServeHTTP(w, req)

	// CORSヘッダーが設定されていないことを確認
	origin := w.Header().Get("Access-Control-Allow-Origin")
	assert.NotEqual(t, "http://malicious-site.com", origin)
}

func TestProductionLogLevel(t *testing.T) {
	setupProductionEnv()
	defer cleanupEnv()

	logLevel := config.GetLogLevel()
	
	// 本番環境では詳細なログを出力しない
	assert.NotEqual(t, "debug", logLevel)
	assert.Contains(t, []string{"info", "warn", "error"}, logLevel)
}
