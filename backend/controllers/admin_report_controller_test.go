package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/domain"

	"github.com/gin-gonic/gin"
)

type fakeAdminReportService struct {
	summary      *domain.AdminSummaryReport
	eventReports []domain.AdminEventReport
	err          error
}

func (service fakeAdminReportService) GetSummary() (*domain.AdminSummaryReport, error) {
	if service.err != nil {
		return nil, service.err
	}

	return service.summary, nil
}

func (service fakeAdminReportService) GetEventReports() ([]domain.AdminEventReport, error) {
	if service.err != nil {
		return nil, service.err
	}

	return service.eventReports, nil
}

func TestAdminReportSummaryReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminReportController(fakeAdminReportService{
		summary: &domain.AdminSummaryReport{TotalEvents: 1},
	})
	router := gin.New()
	router.GET("/admin/reportes/resumen", controller.Summary)

	request := httptest.NewRequest(http.MethodGet, "/admin/reportes/resumen", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAdminReportSummaryErrorReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminReportController(fakeAdminReportService{err: errors.New("boom")})
	router := gin.New()
	router.GET("/admin/reportes/resumen", controller.Summary)

	request := httptest.NewRequest(http.MethodGet, "/admin/reportes/resumen", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
}

func TestAdminReportEventReportsReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminReportController(fakeAdminReportService{
		eventReports: []domain.AdminEventReport{{EventID: 1, ActiveTickets: 10}},
	})
	router := gin.New()
	router.GET("/admin/reportes/eventos", controller.EventReports)

	request := httptest.NewRequest(http.MethodGet, "/admin/reportes/eventos", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAdminReportEventReportsEmptyArray(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminReportController(fakeAdminReportService{
		eventReports: []domain.AdminEventReport{},
	})
	router := gin.New()
	router.GET("/admin/reportes/eventos", controller.EventReports)

	request := httptest.NewRequest(http.MethodGet, "/admin/reportes/eventos", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), `"data":[]`) {
		t.Fatalf("expected empty data array, got %s", response.Body.String())
	}
}

func TestAdminReportEventReportsErrorReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminReportController(fakeAdminReportService{err: errors.New("boom")})
	router := gin.New()
	router.GET("/admin/reportes/eventos", controller.EventReports)

	request := httptest.NewRequest(http.MethodGet, "/admin/reportes/eventos", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
}
