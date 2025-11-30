package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go-app-base/config"
)

func TestIsProduction(t *testing.T) {
	os.Setenv("ENV_MODE", "production")
	defer os.Unsetenv("ENV_MODE")
	
	assert.True(t, config.IsProduction())
	assert.False(t, config.IsDevelopment())
}

func TestIsDevelopment(t *testing.T) {
	os.Setenv("ENV_MODE", "development")
	defer os.Unsetenv("ENV_MODE")
	
	assert.True(t, config.IsDevelopment())
	assert.False(t, config.IsProduction())
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")
	
	assert.Equal(t, "test_value", config.GetEnv("TEST_VAR", "default"))
	assert.Equal(t, "default", config.GetEnv("NON_EXISTENT", "default"))
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	
	assert.Equal(t, 42, config.GetEnvInt("TEST_INT", 0))
	assert.Equal(t, 10, config.GetEnvInt("NON_EXISTENT", 10))
}

func TestGetEnvBool(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")
	
	assert.True(t, config.GetEnvBool("TEST_BOOL", false))
	assert.False(t, config.GetEnvBool("NON_EXISTENT", false))
}

func TestGetEnvDuration(t *testing.T) {
	os.Setenv("TEST_DURATION", "5m")
	defer os.Unsetenv("TEST_DURATION")
	
	assert.Equal(t, 5*time.Minute, config.GetEnvDuration("TEST_DURATION", time.Hour))
	assert.Equal(t, time.Hour, config.GetEnvDuration("NON_EXISTENT", time.Hour))
}

func TestGetLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_LEVEL")
	
	assert.Equal(t, "debug", config.GetLogLevel())
}

func TestGetAllowedOrigins(t *testing.T) {
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	defer os.Unsetenv("ALLOWED_ORIGINS")
	
	origins := config.GetAllowedOrigins()
	assert.Len(t, origins, 2)
	assert.Contains(t, origins, "http://localhost:3000")
	assert.Contains(t, origins, "http://localhost:5173")
}
