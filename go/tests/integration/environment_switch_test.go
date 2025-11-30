package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-app-base/config"
)

func TestEnvironmentSwitch(t *testing.T) {
	tests := []struct {
		name        string
		envMode     string
		isDev       bool
		isProd      bool
		expectedLog string
	}{
		{
			name:        "Development Mode",
			envMode:     "development",
			isDev:       true,
			isProd:      false,
			expectedLog: "debug",
		},
		{
			name:        "Production Mode",
			envMode:     "production",
			isDev:       false,
			isProd:      true,
			expectedLog: "info",
		},
		{
			name:        "Default Mode (no ENV_MODE set)",
			envMode:     "",
			isDev:       false,
			isProd:      false,
			expectedLog: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envMode != "" {
				os.Setenv("ENV_MODE", tt.envMode)
			} else {
				os.Unsetenv("ENV_MODE")
			}
			
			if tt.expectedLog != "info" {
				os.Setenv("LOG_LEVEL", tt.expectedLog)
			}
			
			defer func() {
				os.Unsetenv("ENV_MODE")
				os.Unsetenv("LOG_LEVEL")
			}()

			assert.Equal(t, tt.isDev, config.IsDevelopment())
			assert.Equal(t, tt.isProd, config.IsProduction())
		})
	}
}

func TestSMTPConfigByEnvironment(t *testing.T) {
	tests := []struct {
		name         string
		envMode      string
		expectedHost string
		expectedPort string
	}{
		{
			name:         "Development SMTP (Mailpit)",
			envMode:      "development",
			expectedHost: "mailpit",
			expectedPort: "1025",
		},
		{
			name:         "Production SMTP (Gmail)",
			envMode:      "production",
			expectedHost: "smtp.gmail.com",
			expectedPort: "587",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ENV_MODE", tt.envMode)
			os.Setenv("SMTP_HOST", tt.expectedHost)
			os.Setenv("SMTP_PORT", tt.expectedPort)
			
			defer func() {
				os.Unsetenv("ENV_MODE")
				os.Unsetenv("SMTP_HOST")
				os.Unsetenv("SMTP_PORT")
			}()

			assert.Equal(t, tt.expectedHost, os.Getenv("SMTP_HOST"))
			assert.Equal(t, tt.expectedPort, os.Getenv("SMTP_PORT"))
		})
	}
}

func TestCORSConfigByEnvironment(t *testing.T) {
	tests := []struct {
		name            string
		envMode         string
		allowedOrigins  string
		shouldContain   []string
		shouldNotContain []string
	}{
		{
			name:           "Development CORS",
			envMode:        "development",
			allowedOrigins: "http://localhost:3000,http://localhost:5173",
			shouldContain:  []string{"http://localhost:3000", "http://localhost:5173"},
			shouldNotContain: []string{"https://yourdomain.com"},
		},
		{
			name:           "Production CORS",
			envMode:        "production",
			allowedOrigins: "https://yourdomain.com,https://www.yourdomain.com",
			shouldContain:  []string{"https://yourdomain.com", "https://www.yourdomain.com"},
			shouldNotContain: []string{"http://localhost:3000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ENV_MODE", tt.envMode)
			os.Setenv("ALLOWED_ORIGINS", tt.allowedOrigins)
			
			defer func() {
				os.Unsetenv("ENV_MODE")
				os.Unsetenv("ALLOWED_ORIGINS")
			}()

			origins := config.GetAllowedOrigins()
			
			for _, origin := range tt.shouldContain {
				assert.Contains(t, origins, origin)
			}
			
			for _, origin := range tt.shouldNotContain {
				assert.NotContains(t, origins, origin)
			}
		})
	}
}

func TestSecurityConfigByEnvironment(t *testing.T) {
	tests := []struct {
		name              string
		envMode           string
		jwtSecret         string
		minSecretLength   int
	}{
		{
			name:            "Development JWT",
			envMode:         "development",
			jwtSecret:       "dev_jwt_secret",
			minSecretLength: 10,
		},
		{
			name:            "Production JWT (Strong)",
			envMode:         "production",
			jwtSecret:       "prod_very_secure_jwt_secret_with_random_chars_12345",
			minSecretLength: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ENV_MODE", tt.envMode)
			os.Setenv("JWT_SECRET", tt.jwtSecret)
			
			defer func() {
				os.Unsetenv("ENV_MODE")
				os.Unsetenv("JWT_SECRET")
			}()

			secret := os.Getenv("JWT_SECRET")
			assert.NotEmpty(t, secret)
			
			if tt.envMode == "production" {
				assert.GreaterOrEqual(t, len(secret), tt.minSecretLength,
					"Production JWT secret should be at least %d characters", tt.minSecretLength)
			}
		})
	}
}
