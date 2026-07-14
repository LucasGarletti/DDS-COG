package controllers

import (
	"encoding/json"
	"errors"
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
	registerErr    error
	loginErr       error
	registerInput  services.RegisterInput
	loginInput     services.LoginInput
}

func (service *fakeAuthService) Register(input services.RegisterInput) (*domain.User, error) {
	service.registerInput = input
	if service.registerErr != nil {
		return nil, service.registerErr
	}
	return service.registeredUser, nil
}

func (service *fakeAuthService) Login(input services.LoginInput) (*services.LoginOutput, error) {
	service.loginInput = input
	if service.loginErr != nil {
		return nil, service.loginErr
	}
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

	controller := NewAuthController(&fakeAuthService{
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

	controller := NewAuthController(&fakeAuthService{
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

func TestRegisterTrimsAndNormalizesEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &fakeAuthService{
		registeredUser: &domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@mail.com",
			Role:  domain.UserRoleClient,
		},
	}
	controller := NewAuthController(service)
	router := gin.New()
	router.POST("/auth/register", controller.Register)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{"name":" Test User ","email":"TEST@MAIL.COM","password":" password "}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}
	if service.registerInput.Email != "test@mail.com" {
		t.Fatalf("expected normalized email, got %q", service.registerInput.Email)
	}
}

func TestRegisterWithDuplicateEmailReturnsConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(&fakeAuthService{registerErr: services.ErrEmailAlreadyExists})
	router := gin.New()
	router.POST("/auth/register", controller.Register)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{"name":"Test User","email":"test@mail.com","password":"password"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", response.Code)
	}
}

func TestRegisterWithServiceErrorReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(&fakeAuthService{registerErr: errors.New("service error")})
	router := gin.New()
	router.POST("/auth/register", controller.Register)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{"name":"Test User","email":"test@mail.com","password":"password"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
}

func TestLoginWithInvalidCredentialsReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(&fakeAuthService{loginErr: services.ErrInvalidCredentials})
	router := gin.New()
	router.POST("/auth/login", controller.Login)

	request := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		strings.NewReader(`{"email":"missing@mail.com","password":"wrong"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestLoginWithServiceErrorReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAuthController(&fakeAuthService{loginErr: errors.New("service error")})
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

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
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
