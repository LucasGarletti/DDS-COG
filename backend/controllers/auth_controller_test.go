package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterWithInvalidBodyReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(nil)
	router := gin.New()
	router.POST("/auth/register", controller.Register)

	request := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":"invalid"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestLoginWithInvalidBodyReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(nil)
	router := gin.New()
	router.POST("/auth/login", controller.Login)

	request := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"invalid"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestMeWithUserInContextReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(nil)
	router := gin.New()
	router.GET("/auth/me", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("name", "Test User")
		c.Set("email", "test@mail.com")
		controller.Me(c)
	})

	request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}
