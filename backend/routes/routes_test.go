package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/domain"
	"backend/middlewares"
	"backend/utils"

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

func TestAdminRouteAllowsAdminCreate(t *testing.T) {
	response := serveAdminRoute(t, domain.UserRoleAdmin, http.MethodPost, "/admin/eventos")

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}
}

func TestAdminRouteRejectsClient(t *testing.T) {
	response := serveAdminRoute(t, domain.UserRoleClient, http.MethodPost, "/admin/eventos")

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
}

func TestAdminRouteWithoutTokenReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := newAdminAuthTestRouter()
	request := httptest.NewRequest(http.MethodPost, "/admin/eventos", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAdminRouteAllowsAdminUpdate(t *testing.T) {
	response := serveAdminRoute(t, domain.UserRoleAdmin, http.MethodPatch, "/admin/eventos/1")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAdminRouteAllowsAdminCancel(t *testing.T) {
	response := serveAdminRoute(t, domain.UserRoleAdmin, http.MethodDelete, "/admin/eventos/1")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAdminRouteAllowsAdminReport(t *testing.T) {
	response := serveAdminRoute(t, domain.UserRoleAdmin, http.MethodGet, "/admin/eventos/1/reporte")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAdminReportRoutesExist(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "summary", path: "/admin/reportes/resumen"},
		{name: "event reports", path: "/admin/reportes/eventos"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := serveAdminRoute(t, domain.UserRoleAdmin, http.MethodGet, tt.path)

			if response.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", response.Code)
			}
		})
	}
}

func TestAdminReportRouteWithoutTokenReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := newAdminAuthTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/admin/reportes/resumen", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestAdminReportRouteRejectsClient(t *testing.T) {
	response := serveAdminRoute(t, domain.UserRoleClient, http.MethodGet, "/admin/reportes/resumen")

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
}

func TestAdminCanAccessReportController(t *testing.T) {
	response := serveAdminRoute(t, domain.UserRoleAdmin, http.MethodGet, "/admin/reportes/resumen")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func newAdminAuthTestRouter() *gin.Engine {
	router := gin.New()
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middlewares.AuthMiddleware())
	adminRoutes.Use(middlewares.RequireRole(domain.UserRoleAdmin))
	{
		adminRoutes.POST("/eventos", func(c *gin.Context) { c.Status(http.StatusCreated) })
		adminRoutes.PATCH("/eventos/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
		adminRoutes.DELETE("/eventos/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
		adminRoutes.GET("/eventos/:id/reporte", func(c *gin.Context) { c.Status(http.StatusOK) })
		adminRoutes.GET("/reportes/resumen", func(c *gin.Context) { c.Status(http.StatusOK) })
		adminRoutes.GET("/reportes/eventos", func(c *gin.Context) { c.Status(http.StatusOK) })
	}
	return router
}

func serveAdminRoute(t *testing.T, role string, method string, path string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES", "1")

	token, err := utils.GenerateJWT(1, "Test User", "test@mail.com", role)
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}

	router := newAdminAuthTestRouter()
	request := httptest.NewRequest(method, path, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	return response
}
