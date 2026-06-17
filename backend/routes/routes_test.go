package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPingReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := SetupRouter(nil)
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestProtectedEndpointWithoutTokenReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := SetupRouter(nil)
	request := httptest.NewRequest(http.MethodGet, "/mis-entradas", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestTicketProtectedEndpointsWithoutTokenReturnUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "purchase", method: http.MethodPost, path: "/entradas/comprar/1"},
		{name: "my tickets", method: http.MethodGet, path: "/mis-entradas"},
		{name: "cancel", method: http.MethodPatch, path: "/entradas/1/cancelar"},
		{name: "transfer", method: http.MethodPatch, path: "/entradas/1/transferir"},
	}

	router := SetupRouter(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusUnauthorized {
				t.Fatalf("expected status 401, got %d", response.Code)
			}
		})
	}
}
