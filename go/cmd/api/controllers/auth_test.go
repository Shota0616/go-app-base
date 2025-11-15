package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"go-app-base/auth"
	"go-app-base/config"
	"go-app-base/cmd/api/controllers"
	"go-app-base/middleware"
	"go-app-base/models"
	"gorm.io/gorm"
	"path/filepath"
	"golang.org/x/text/language"
)

var (
	router *gin.Engine
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
		public.POST("/verify", controllers.Verify)
		public.POST("/resend-verification-code", controllers.ResendVerificationCode)
		public.POST("/request-password-reset", controllers.RequestPasswordReset)
		public.POST("/reset-password", controllers.ResetPassword)
		public.POST("/refresh-token", controllers.RefreshToken)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/getuser", controllers.GetUser)
		protected.PUT("/user/username", controllers.UpdateUsername)
		protected.PUT("/user/email", controllers.UpdateEmail)
		protected.PUT("/user/password", controllers.UpdatePassword)
		protected.DELETE("/user", controllers.DeleteUser)
	}
	return r
}

func TestMain(m *testing.M) {
	// Setup environment variables for testing
	os.Setenv("JWT_SECRET", "test_jwt_secret")
	os.Setenv("APP_URL", "http://localhost:3000")
	os.Setenv("EMAIL_FROM", "noreply@example.com")
	os.Setenv("SMTP_HOST", "localhost")
	os.Setenv("SMTP_PORT", "1025")
	os.Setenv("APP_ROOT", "/usr/src") // Set APP_ROOT for i18n to find locale files
	os.Setenv("APP_LANG", "en") // Explicitly set language for i18n in tests

	// DBとRedisを本番同様に初期化
	config.ConnectDatabase()
	config.MigrateDatabase()
	config.ConnectRedis()

	// Initialize i18n for controllers
	// config.InitI18n() // Commented out

	// Directly initialize config.Localizer for testing
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	enPath := filepath.Join(os.Getenv("APP_ROOT"), "locales", "en.json")
	jaPath := filepath.Join(os.Getenv("APP_ROOT"), "locales", "ja.json")
	if _, err := bundle.LoadMessageFile(enPath); err != nil {
		panic("Failed to load en.json in test: " + err.Error())
	}
	if _, err := bundle.LoadMessageFile(jaPath); err != nil {
		panic("Failed to load ja.json in test: " + err.Error())
	}
	config.Localizer = i18n.NewLocalizer(bundle, "en")

	router = setupRouter()

	exitCode := m.Run()

	// Teardown
	sqlDB, _ := config.DB.DB()
	sqlDB.Close()
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("APP_URL")
	os.Unsetenv("EMAIL_FROM")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("APP_ROOT")
	os.Unsetenv("APP_LANG")

	os.Exit(exitCode)
}

// Helper function to create a Gin context and recorder
func getTestContext(method, url string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, url, nil)
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")
	}
	return c, w
}

// Helper to create a user in the test DB
func createUser(email, password, username string, isActive bool) models.User {
	hashedPassword, _ := auth.HashPassword(password)
	user := models.User{
		Email:    email,
		Password: hashedPassword,
		Username: username,
		IsActive: isActive,
	}
	config.DB.Create(&user)
	return user
}

// --- Test Register Function ---
func TestRegisterSuccess(t *testing.T) {
	// Reset DB before each test
	config.DB.Exec("DELETE FROM users")

	body := gin.H{"email": "test@example.com", "password": "password123"}
	c, w := getTestContext(http.MethodPost, "/api/register", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_registration_success"}), response["message"])

	var user models.User
	config.DB.Where("email = ?", "test@example.com").First(&user)
	assert.False(t, user.IsActive)
	assert.NotEmpty(t, user.Username)
}

func TestRegisterDuplicateEmail(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	createUser("test@example.com", "password123", "testuser", false)

	body := gin.H{"email": "test@example.com", "password": "password123"}
	c, w := getTestContext(http.MethodPost, "/api/register", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusConflict, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_already_registered"}), response["error"])
}

func TestRegisterInvalidInput(t *testing.T) {
	config.DB.Exec("DELETE FROM users")

	body := gin.H{"email": "invalid-email", "password": "password123"} // Invalid email format
	c, w := getTestContext(http.MethodPost, "/api/register", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "input_data_invalid"}), response["error"])
}

// --- Test Verify Function ---
func TestVerifySuccess_Registration(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("verify@example.com", "password123", "verifyuser", false)
	verificationCode := "123456"

	body := gin.H{"email": user.Email, "verificationCode": verificationCode}
	c, w := getTestContext(http.MethodPost, "/api/verify", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_verified"}), response["message"])

	var updatedUser models.User
	config.DB.Where("id = ?", user.ID).First(&updatedUser)
	assert.True(t, updatedUser.IsActive)
}

func TestVerifySuccess_EmailChange(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("old@example.com", "password123", "changeuser", true) // Active user
	newEmail := "new@example.com"
	verificationCode := "654321"

	body := gin.H{"email": newEmail, "verificationCode": verificationCode}
	c, w := getTestContext(http.MethodPost, "/api/verify", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_verified"}), response["message"])

	var updatedUser models.User
	config.DB.Where("id = ?", user.ID).First(&updatedUser)
	assert.True(t, updatedUser.IsActive)
	assert.Equal(t, newEmail, updatedUser.Email)
}

func TestVerifyInvalidCode(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("invalidcode@example.com", "password123", "invalidcodeuser", false)
	wrongCode := "999999"

	body := gin.H{"email": user.Email, "verificationCode": wrongCode}
	c, w := getTestContext(http.MethodPost, "/api/verify", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_or_expired_verification_code"}), response["error"])

	var updatedUser models.User
	config.DB.Where("id = ?", user.ID).First(&updatedUser)
	assert.False(t, updatedUser.IsActive) // Should still be inactive
}

func TestVerifyExpiredCode(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("expiredcode@example.com", "password123", "expiredcodeuser", false)

	body := gin.H{"email": user.Email, "verificationCode": "123456"}
	c, w := getTestContext(http.MethodPost, "/api/verify", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_or_expired_verification_code"}), response["error"])

	var updatedUser models.User
	config.DB.Where("id = ?", user.ID).First(&updatedUser)
	assert.False(t, updatedUser.IsActive) // Should still be inactive
}

// --- Test ResendVerificationCode Function ---
func TestResendVerificationCodeSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("resend@example.com", "password123", "resenduser", false)

	body := gin.H{"email": user.Email}
	c, w := getTestContext(http.MethodPost, "/api/resend-verification-code", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "verification_code_resent"}), response["message"])
}

func TestResendVerificationCodeLimitReached(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("resendlimit@example.com", "password123", "resendlimituser", false)

	body := gin.H{"email": user.Email}
	c, w := getTestContext(http.MethodPost, "/api/resend-verification-code", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resend_limit_reached"}), response["error"])
}

// --- Test Login Function ---
func TestLoginSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	password := "password123"
	user := createUser("login@example.com", password, "loginuser", true) // Active user

	body := gin.H{"email": user.Email, "password": password}
	c, w := getTestContext(http.MethodPost, "/api/login", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"])
	assert.NotEmpty(t, response["refreshtoken"])
}

func TestLoginInvalidCredentials(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("invalidlogin@example.com", "password123", "invalidloginuser", true)

	body := gin.H{"email": user.Email, "password": "wrongpassword"}
	c, w := getTestContext(http.MethodPost, "/api/login", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_incorrect"}), response["error"])
}

func TestLoginInactiveAccount(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("inactive@example.com", "password123", "inactiveuser", false) // Inactive user

	body := gin.H{"email": user.Email, "password": "password123"}
	c, w := getTestContext(http.MethodPost, "/api/login", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusSeeOther, w.Code) // StatusSeeOther for inactive account
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "account_not_activated_resend_verification"}), response["error"])
}

// --- Test RequestPasswordReset Function ---
func TestRequestPasswordResetSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("reset@example.com", "password123", "resetuser", true) // Active user

	body := gin.H{"email": user.Email}
	c, w := getTestContext(http.MethodPost, "/api/request-password-reset", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset_link_sent"}), response["message"])
}

func TestRequestPasswordResetInactiveAccount(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("resetinactive@example.com", "password123", "resetinactiveuser", false) // Inactive user

	body := gin.H{"email": user.Email}
	c, w := getTestContext(http.MethodPost, "/api/request-password-reset", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "account_not_activated_resend_verification"}), response["error"])
}

// --- Test ResetPassword Function ---
func TestResetPasswordSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("resetpass@example.com", "oldpassword", "resetpassuser", true)
	token, _ := auth.GenerateJWT(user.ID) // Generate a valid token

	body := gin.H{"token": token, "newPassword": "newpassword123"}
	c, w := getTestContext(http.MethodPost, "/api/reset-password", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_reset_success"}), response["message"])

	var updatedUser models.User
	config.DB.Where("id = ?", user.ID).First(&updatedUser)
	assert.True(t, auth.CheckPasswordHash("newpassword123", updatedUser.Password))
}

func TestResetPasswordInvalidToken(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	createUser("resetinvalid@example.com", "oldpassword", "resetinvaliduser", true)

	// Mock auth.ValidateJWT to return an error for any token
	originalValidateJWT := auth.ValidateJWT
	auth.ValidateJWT = func(tokenStr string) (*auth.Claims, error) {
		return nil, fmt.Errorf("invalid or expired token")
	}
	defer func() { auth.ValidateJWT = originalValidateJWT }()

	body := gin.H{"token": "invalid.token.here", "newPassword": "newpassword123"}
	c, w := getTestContext(http.MethodPost, "/api/reset-password", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_or_expired_token"}), response["error"])
}

// --- Test GetUser Function ---
func TestGetUserSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("getuser@example.com", "password123", "getuser", true)
	token, _ := auth.GenerateJWT(user.ID)

	c, w := getTestContext(http.MethodGet, "/api/getuser", nil)
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(user.ID), response["id"]) // JSON unmarshals numbers to float64
	assert.Equal(t, user.Username, response["username"])
	assert.Equal(t, user.Email, response["email"])
	assert.Equal(t, user.IsActive, response["active"])
}

func TestGetUserUnauthorized(t *testing.T) {
	config.DB.Exec("DELETE FROM users")

	c, w := getTestContext(http.MethodGet, "/api/getuser", nil)
	// No Authorization header
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "authorization_token_not_provided"}), response["error"])
}

// --- Test UpdateUsername Function ---
func TestUpdateUsernameSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("updateusername@example.com", "password123", "oldusername", true)
	token, _ := auth.GenerateJWT(user.ID)
	newUsername := "newusername"

	body := gin.H{"username": newUsername}
	c, w := getTestContext(http.MethodPut, "/api/user/username", body)
	c.Set("id", user.ID) // Mock AuthRequired middleware
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "username_updated_successfully"}), response["message"])

	var updatedUser models.User
	config.DB.Where("id = ?", user.ID).First(&updatedUser)
	assert.Equal(t, newUsername, updatedUser.Username)
}

func TestUpdateUsernameDuplicate(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user1 := createUser("updateusername1@example.com", "password123", "existingusername", true)
	user2 := createUser("updateusername2@example.com", "password123", "anotheruser", true)
	token, _ := auth.GenerateJWT(user2.ID) // User2 tries to update to user1's username

	body := gin.H{"username": user1.Username}
	c, w := getTestContext(http.MethodPut, "/api/user/username", body)
	c.Set("id", user2.ID) // Mock AuthRequired middleware
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusConflict, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "username_already_registered"}), response["error"])
}

// --- Test UpdateEmail Function ---
func TestUpdateEmailSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("oldemail@example.com", "password123", "emailuser", true)
	token, _ := auth.GenerateJWT(user.ID)
	newEmail := "newemail@example.com"

	body := gin.H{"newEmail": newEmail}
	c, w := getTestContext(http.MethodPut, "/api/user/email", body)
	c.Set("id", user.ID) // Mock AuthRequired middleware
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_updated_successfully_please_verify"}), response["message"])

	// Verify that user's email in DB is NOT changed yet
	var originalUser models.User
	config.DB.Where("id = ?", user.ID).First(&originalUser)
	assert.Equal(t, user.Email, originalUser.Email) // Should still be old email
}

func TestUpdateEmailDuplicateNewEmail(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user1 := createUser("existing@example.com", "password123", "user1", true)
	user2 := createUser("changer@example.com", "password123", "user2", true)
	token, _ := auth.GenerateJWT(user2.ID)

	body := gin.H{"newEmail": user1.Email} // User2 tries to change to user1's email
	c, w := getTestContext(http.MethodPut, "/api/user/email", body)
	c.Set("id", user2.ID)
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusConflict, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_already_registered"}), response["error"])
}

func TestUpdateEmailSameEmail(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("same@example.com", "password123", "sameuser", true)
	token, _ := auth.GenerateJWT(user.ID)

	body := gin.H{"newEmail": user.Email}
	c, w := getTestContext(http.MethodPut, "/api/user/email", body)
	c.Set("id", user.ID)
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email_already_current"}), response["message"])
}

// --- Test UpdatePassword Function ---
func TestUpdatePasswordSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	oldPassword := "oldpassword123"
	newPassword := "newpassword456"
	user := createUser("updatepass@example.com", oldPassword, "updatepassuser", true)
	token, _ := auth.GenerateJWT(user.ID)

	body := gin.H{"currentPassword": oldPassword, "newPassword": newPassword}
	c, w := getTestContext(http.MethodPut, "/api/user/password", body)
	c.Set("id", user.ID)
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "password_updated_successfully"}), response["message"])

	var updatedUser models.User
	config.DB.Where("id = ?", user.ID).First(&updatedUser)
	assert.True(t, auth.CheckPasswordHash(newPassword, updatedUser.Password))
}

func TestUpdatePasswordIncorrectCurrent(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	oldPassword := "oldpassword123"
	newPassword := "newpassword456"
	user := createUser("wrongpass@example.com", oldPassword, "wrongpassuser", true)
	token, _ := auth.GenerateJWT(user.ID)

	body := gin.H{"currentPassword": "wrongcurrentpassword", "newPassword": newPassword}
	c, w := getTestContext(http.MethodPut, "/api/user/password", body)
	c.Set("id", user.ID)
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "incorrect_current_password"}), response["error"])
}

// --- Test DeleteUser Function ---
func TestDeleteUserSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("delete@example.com", "password123", "deleteuser", true)
	token, _ := auth.GenerateJWT(user.ID)

	c, w := getTestContext(http.MethodDelete, "/api/user", nil)
	c.Set("id", user.ID)
	c.Request.Header.Set("Authorization", token)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "user_deleted_successfully"}), response["message"])

	var deletedUser models.User
	err := config.DB.Where("id = ?", user.ID).First(&deletedUser).Error
	assert.Error(t, err) // User should not be found
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDeleteUserUnauthorized(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	createUser("nodelete@example.com", "password123", "nodeleteuser", true)

	c, w := getTestContext(http.MethodDelete, "/api/user", nil)
	// No Authorization header
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "authorization_token_not_provided"}), response["error"])
}

// --- Test RefreshToken Function ---
func TestRefreshTokenSuccess(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	user := createUser("refresh@example.com", "password123", "refreshuser", true)
	refreshToken, _ := auth.GenerateRefreshToken(user.ID)

	// Mock auth.ValidateJWT to return valid claims
	originalValidateJWT := auth.ValidateJWT
	auth.ValidateJWT = func(tokenStr string) (*auth.Claims, error) {
		return &auth.Claims{ID: user.ID}, nil
	}
	defer func() { auth.ValidateJWT = originalValidateJWT }()

	body := gin.H{"refresh_token": refreshToken}
	c, w := getTestContext(http.MethodPost, "/api/refresh-token", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"])
}

func TestRefreshTokenInvalid(t *testing.T) {
	config.DB.Exec("DELETE FROM users")
	createUser("refreshinvalid@example.com", "password123", "refreshinvaliduser", true)

	// Mock auth.ValidateJWT to return an error for any token
	originalValidateJWT := auth.ValidateJWT
	auth.ValidateJWT = func(tokenStr string) (*auth.Claims, error) {
		return nil, fmt.Errorf("invalid refresh token")
	}
	defer func() { auth.ValidateJWT = originalValidateJWT }()

	body := gin.H{"refresh_token": "invalid.refresh.token"}
	c, w := getTestContext(http.MethodPost, "/api/refresh-token", body)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, config.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_refresh_token"}), response["error"])
}
