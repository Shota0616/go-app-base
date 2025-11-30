package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-app-base/config"
	"go-app-base/routes"
)

func setupDevelopmentEnv() {
	os.Setenv("ENV_MODE", "development")
	os.Setenv("JWT_SECRET", "test_jwt_secret")
	os.Setenv("JWT_REFRESH_SECRET", "test_jwt_refresh_secret")
	os.Setenv("APP_URL", "http://localhost:8000")
	os.Setenv("EMAIL_FROM", "noreply@example.com")
	os.Setenv("SMTP_HOST", "mailpit")
	os.Setenv("SMTP_PORT", "1025")
	os.Setenv("REDIS_HOST", "redis")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("APP_ROOT", "/usr/src")
	os.Setenv("APP_LANG", "en")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	os.Setenv("LOG_LEVEL", "debug")
}

func TestDevelopmentEnvironment(t *testing.T) {
	setupDevelopmentEnv()
	defer cleanupEnv()

	assert.True(t, config.IsDevelopment())
	assert.False(t, config.IsProduction())
	assert.Equal(t, "debug", config.GetLogLevel())
}

func TestDevelopmentCORS(t *testing.T) {
	setupDevelopmentEnv()
	defer cleanupEnv()

	origins := config.GetAllowedOrigins()
	assert.Contains(t, origins, "http://localhost:3000")
	assert.Contains(t, origins, "http://localhost:5173")
}

func TestDevelopmentMailpitConfig(t *testing.T) {
	setupDevelopmentEnv()
	defer cleanupEnv()

	assert.Equal(t, "mailpit", os.Getenv("SMTP_HOST"))
	assert.Equal(t, "1025", os.Getenv("SMTP_PORT"))
}

func TestDevelopmentPingEndpoint(t *testing.T) {
	setupDevelopmentEnv()
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

func cleanupEnv() {
	os.Unsetenv("ENV_MODE")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_REFRESH_SECRET")
	os.Unsetenv("APP_URL")
	os.Unsetenv("EMAIL_FROM")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("APP_ROOT")
	os.Unsetenv("APP_LANG")
	os.Unsetenv("ALLOWED_ORIGINS")
	os.Unsetenv("LOG_LEVEL")
}
