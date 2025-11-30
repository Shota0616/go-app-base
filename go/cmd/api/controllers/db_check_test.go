package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-app-base/cmd/api/controllers"
	"go-app-base/config"
)

func TestDBCheckSuccess(t *testing.T) {
	os.Setenv("APP_ROOT", "/usr/src")
	config.ConnectDatabase()
	defer func() {
		sqlDB, _ := config.DB.DB()
		sqlDB.Close()
	}()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/db-check", controllers.DBCheck)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/db-check", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Database connection is healthy", response["message"])
}
