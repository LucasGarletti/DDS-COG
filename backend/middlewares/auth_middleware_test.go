package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/utils"

	"github.com/gin-gonic/gin"
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

	token, err := utils.GenerateJWT(1, "Test User", "test@mail.com")
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
