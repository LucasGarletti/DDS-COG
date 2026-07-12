package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/domain"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func TestRequireRoleAllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protected", func(c *gin.Context) {
		c.Set("role", domain.UserRoleAdmin)
	}, RequireRole(domain.UserRoleAdmin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestRequireRoleRejectsClient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protected", func(c *gin.Context) {
		c.Set("role", domain.UserRoleClient)
	}, RequireRole(domain.UserRoleAdmin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}

	assertJSONError(t, response, "insufficient permissions")
}

func TestRequireRoleRejectsMissingRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protected", RequireRole(domain.UserRoleAdmin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}

	assertJSONError(t, response, "authentication required")
}

func TestRequireRoleRejectsInvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protected", func(c *gin.Context) {
		c.Set("role", "superadmin")
	}, RequireRole(domain.UserRoleAdmin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}

	assertJSONError(t, response, "insufficient permissions")
}

func TestRequireRoleRejectsNonStringRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/protected", func(c *gin.Context) {
		c.Set("role", 123)
	}, RequireRole(domain.UserRoleAdmin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}

	assertJSONError(t, response, "insufficient permissions")
}

func TestRequireRoleEndpointAccessibleForAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES", "1")

	token, err := utils.GenerateJWT(1, "Admin User", "admin@mail.com", domain.UserRoleAdmin)
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware())
	adminGroup := router.Group("/admin")
	adminGroup.Use(RequireRole(domain.UserRoleAdmin))
	adminGroup.GET("/dashboard", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestRequireRoleEndpointBlockedForClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES", "1")

	token, err := utils.GenerateJWT(1, "Client User", "client@mail.com", domain.UserRoleClient)
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware())
	adminGroup := router.Group("/admin")
	adminGroup.Use(RequireRole(domain.UserRoleAdmin))
	adminGroup.GET("/dashboard", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}

	assertJSONError(t, response, "insufficient permissions")
}

func assertJSONError(t *testing.T, response *httptest.ResponseRecorder, expectedError string) {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if body["success"] != false {
		t.Fatalf("expected success false, got %v", body["success"])
	}

	if body["error"] != expectedError {
		t.Fatalf("expected error %q, got %v", expectedError, body["error"])
	}
}
