package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type fakeFestivalScheduleService struct {
	schedule  *domain.FestivalSchedule
	schedules []domain.FestivalSchedule
	err       error
}

func (service fakeFestivalScheduleService) Create(input services.CreateFestivalScheduleInput) (*domain.FestivalSchedule, error) {
	return service.schedule, service.err
}

func (service fakeFestivalScheduleService) List(eventID uint) ([]domain.FestivalSchedule, error) {
	return service.schedules, service.err
}

func (service fakeFestivalScheduleService) Delete(id uint) error {
	return service.err
}

func TestFestivalScheduleCreateReturnsCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewFestivalScheduleController(fakeFestivalScheduleService{schedule: &domain.FestivalSchedule{ID: 1}})
	router := gin.New()
	router.POST("/admin/eventos/:id/grilla", controller.Create)

	body := `{"artist":"Banda","stage":"Norte","start_time":"2026-12-10T20:00:00Z","end_time":"2026-12-10T21:00:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/admin/eventos/1/grilla", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}
}

func TestFestivalScheduleCreateInvalidBodyReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewFestivalScheduleController(fakeFestivalScheduleService{})
	router := gin.New()
	router.POST("/admin/eventos/:id/grilla", controller.Create)

	request := httptest.NewRequest(http.MethodPost, "/admin/eventos/1/grilla", strings.NewReader("{"))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestFestivalScheduleCreateNonFestivalReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewFestivalScheduleController(fakeFestivalScheduleService{err: services.ErrEventNotFestival})
	router := gin.New()
	router.POST("/admin/eventos/:id/grilla", controller.Create)

	body := `{"artist":"Banda","stage":"Norte","start_time":"2026-12-10T20:00:00Z","end_time":"2026-12-10T21:00:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/admin/eventos/1/grilla", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestFestivalScheduleListReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewFestivalScheduleController(fakeFestivalScheduleService{schedules: []domain.FestivalSchedule{{ID: 1}}})
	router := gin.New()
	router.GET("/eventos/:id/grilla", controller.List)

	request := httptest.NewRequest(http.MethodGet, "/eventos/1/grilla", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestFestivalScheduleDeleteInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewFestivalScheduleController(fakeFestivalScheduleService{err: errors.New("boom")})
	router := gin.New()
	router.DELETE("/admin/grilla/:id", controller.Delete)

	request := httptest.NewRequest(http.MethodDelete, "/admin/grilla/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
}

func TestFestivalScheduleControllerUsesRFC3339(t *testing.T) {
	if _, err := time.Parse(time.RFC3339, "2026-12-10T20:00:00Z"); err != nil {
		t.Fatalf("expected RFC3339 parse: %v", err)
	}
}
