package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type fakeAuthService struct {
	registeredUser *domain.User
	loginOutput    *services.LoginOutput
}

func (service fakeAuthService) Register(input services.RegisterInput) (*domain.User, error) {
	return service.registeredUser, nil
}

func (service fakeAuthService) Login(input services.LoginInput) (*services.LoginOutput, error) {
	return service.loginOutput, nil
}

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

func TestRegisterReturnsClientRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(fakeAuthService{
		registeredUser: &domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@mail.com",
			Role:  domain.UserRoleClient,
		},
	})
	router := gin.New()
	router.POST("/auth/register", controller.Register)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{"name":"Test User","email":"test@mail.com","password":"password","role":"admin"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	data := body["data"].(map[string]any)
	if data["role"] != domain.UserRoleClient {
		t.Fatalf("expected role client, got %v", data["role"])
	}
}

func TestLoginReturnsRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(fakeAuthService{
		loginOutput: &services.LoginOutput{
			Token: "token",
			User: &domain.User{
				ID:    1,
				Name:  "Test User",
				Email: "test@mail.com",
				Role:  domain.UserRoleClient,
			},
		},
	})
	router := gin.New()
	router.POST("/auth/login", controller.Login)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		strings.NewReader(`{"email":"test@mail.com","password":"password"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	data := body["data"].(map[string]any)
	if data["role"] != domain.UserRoleClient {
		t.Fatalf("expected role client, got %v", data["role"])
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
		c.Set("role", domain.UserRoleClient)
		controller.Me(c)
	})

	request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	data := body["data"].(map[string]any)
	if data["role"] != domain.UserRoleClient {
		t.Fatalf("expected role client, got %v", data["role"])
	}
}
