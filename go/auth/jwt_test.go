package auth_test

import (
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"go-app-base/auth"
)

func TestGenerateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	userID := uint(123)
	token, err := auth.GenerateJWT(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	userID := uint(456)
	token, _ := auth.GenerateJWT(userID)

	claims, err := auth.ValidateJWT(token)

	assert.NoError(t, err)
	assert.Equal(t, userID, claims.ID)
}

func TestValidateJWTInvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	invalidToken := "invalid.token.string"
	_, err := auth.ValidateJWT(invalidToken)

	assert.Error(t, err)
}

func TestValidateJWTExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	expirationTime := time.Now().Add(-1 * time.Hour)
	claims := &auth.Claims{
		ID: 789,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test_secret"))

	_, err := auth.ValidateJWT(tokenString)
	assert.Error(t, err)
}

func TestHashPassword(t *testing.T) {
	password := "mypassword123"
	hashedPassword, err := auth.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mypassword123"
	hashedPassword, _ := auth.HashPassword(password)

	result := auth.CheckPasswordHash(password, hashedPassword)
	assert.True(t, result)

	wrongPassword := "wrongpassword"
	result = auth.CheckPasswordHash(wrongPassword, hashedPassword)
	assert.False(t, result)
}
