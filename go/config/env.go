package config

import (
	"os"
	"strconv"
	"time"
)

// IsProduction returns true if running in production mode
func IsProduction() bool {
	return os.Getenv("ENV_MODE") == "production"
}

// IsDevelopment returns true if running in development mode
func IsDevelopment() bool {
	return os.Getenv("ENV_MODE") == "development"
}

// GetEnv gets environment variable with default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt gets environment variable as integer with default value
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool gets environment variable as boolean with default value
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetEnvDuration gets environment variable as duration with default value
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetLogLevel returns the configured log level
func GetLogLevel() string {
	return GetEnv("LOG_LEVEL", "info")
}

// GetAllowedOrigins returns the list of allowed CORS origins
func GetAllowedOrigins() []string {
	origins := GetEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	if origins == "" {
		return []string{}
	}
	
	var result []string
	for i := 0; i < len(origins); {
		end := i
		for end < len(origins) && origins[end] != ',' {
			end++
		}
		if i < end {
			result = append(result, origins[i:end])
		}
		i = end + 1
	}
	return result
}
