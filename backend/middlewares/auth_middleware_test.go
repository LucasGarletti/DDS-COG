package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/domain"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddlewareWithoutAuthorizationReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAuthMiddlewareWithInvalidAuthorizationReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")

	router := gin.New()
	router.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer invalid-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAuthMiddlewareWithValidAuthorizationContinues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES", "1")

	token, err := utils.GenerateJWT(1, "Test User", "test@mail.com", domain.UserRoleClient)
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}

	router := gin.New()
	router.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			t.Fatal("expected user_id in context")
		}

		if userID != uint(1) {
			t.Fatalf("expected user_id 1, got %v", userID)
		}

		role, exists := c.Get("role")
		if !exists {
			t.Fatal("expected role in context")
		}

		if role != domain.UserRoleClient {
			t.Fatalf("expected role client, got %v", role)
		}

		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAuthMiddlewareWithInvalidRoleReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")

	claims := utils.JWTClaims{
		ID:    1,
		Name:  "Test User",
		Email: "test@mail.com",
		Role:  "superadmin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("could not sign token: %v", err)
	}

	router := gin.New()
	router.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}
