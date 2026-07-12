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
	"gorm.io/gorm"
)

type fakeAdminEventService struct {
	event  *domain.Event
	report *services.EventReport
	err    error
}

func (service fakeAdminEventService) CreateEvent(input services.CreateEventInput) (*domain.Event, error) {
	return service.event, service.err
}

func (service fakeAdminEventService) UpdateEvent(input services.UpdateEventInput) (*domain.Event, error) {
	return service.event, service.err
}

func (service fakeAdminEventService) CancelEvent(id uint) (*domain.Event, error) {
	return service.event, service.err
}

func (service fakeAdminEventService) GetEventReport(id uint) (*services.EventReport, error) {
	return service.report, service.err
}

type trackingAdminEventService struct {
	createInput services.CreateEventInput
	updateInput services.UpdateEventInput
	event       *domain.Event
}

func (service *trackingAdminEventService) CreateEvent(input services.CreateEventInput) (*domain.Event, error) {
	service.createInput = input
	return service.event, nil
}

func (service *trackingAdminEventService) UpdateEvent(input services.UpdateEventInput) (*domain.Event, error) {
	service.updateInput = input
	return service.event, nil
}

func (service *trackingAdminEventService) CancelEvent(id uint) (*domain.Event, error) {
	return service.event, nil
}

func (service *trackingAdminEventService) GetEventReport(id uint) (*services.EventReport, error) {
	return nil, nil
}

func validAdminEventJSON() string {
	return `{"title":"Evento","description":"Descripcion","date":"2026-12-10T20:00:00-03:00","location":"Cordoba","capacity":1000,"price":25000,"image_url":"https://example.com/image.jpg"}`
}

func TestAdminCreateEventReturnsCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{
		event: &domain.Event{ID: 1, Title: "Evento", Status: domain.EventStatusActive},
	})
	router := gin.New()
	router.POST("/admin/eventos", controller.Create)

	request := httptest.NewRequest(http.MethodPost, "/admin/eventos", strings.NewReader(validAdminEventJSON()))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}
}

func TestAdminCreateEventPassesIsFestival(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &trackingAdminEventService{event: &domain.Event{ID: 1, IsFestival: true}}
	controller := NewAdminEventController(service)
	router := gin.New()
	router.POST("/admin/eventos", controller.Create)

	body := `{"title":"Evento","description":"Descripcion","date":"2026-12-10T20:00:00-03:00","location":"Cordoba","capacity":1000,"price":25000,"image_url":"https://example.com/image.jpg","is_festival":true}`
	request := httptest.NewRequest(http.MethodPost, "/admin/eventos", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}

	if !service.createInput.IsFestival {
		t.Fatal("expected is_festival true to be passed to service")
	}
}

func TestAdminCreateEventInvalidJSONReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{})
	router := gin.New()
	router.POST("/admin/eventos", controller.Create)

	request := httptest.NewRequest(http.MethodPost, "/admin/eventos", strings.NewReader("{"))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestAdminCreateEventInvalidValidationReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{err: services.ErrInvalidEventData})
	router := gin.New()
	router.POST("/admin/eventos", controller.Create)

	request := httptest.NewRequest(http.MethodPost, "/admin/eventos", strings.NewReader(validAdminEventJSON()))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestAdminUpdateMissingEventReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{err: gorm.ErrRecordNotFound})
	router := gin.New()
	router.PATCH("/admin/eventos/:id", controller.Update)

	request := httptest.NewRequest(http.MethodPatch, "/admin/eventos/1", strings.NewReader(`{"title":"Nuevo"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestAdminUpdateConflictReturnsConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{err: services.ErrEventCapacityConflict})
	router := gin.New()
	router.PATCH("/admin/eventos/:id", controller.Update)

	request := httptest.NewRequest(http.MethodPatch, "/admin/eventos/1", strings.NewReader(`{"capacity":1}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", response.Code)
	}
}

func TestAdminCreateInternalErrorReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{err: errors.New("internal")})
	router := gin.New()
	router.POST("/admin/eventos", controller.Create)

	request := httptest.NewRequest(http.MethodPost, "/admin/eventos", strings.NewReader(validAdminEventJSON()))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
}

func TestAdminUpdateEventReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{
		event: &domain.Event{ID: 1, Title: "Nuevo"},
	})
	router := gin.New()
	router.PATCH("/admin/eventos/:id", controller.Update)

	request := httptest.NewRequest(http.MethodPatch, "/admin/eventos/1", strings.NewReader(`{"title":"Nuevo"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAdminUpdateEventPassesIsFestival(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &trackingAdminEventService{event: &domain.Event{ID: 1, IsFestival: true}}
	controller := NewAdminEventController(service)
	router := gin.New()
	router.PATCH("/admin/eventos/:id", controller.Update)

	request := httptest.NewRequest(http.MethodPatch, "/admin/eventos/1", strings.NewReader(`{"is_festival":true}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	if service.updateInput.IsFestival == nil || !*service.updateInput.IsFestival {
		t.Fatal("expected is_festival true pointer to be passed to service")
	}
}

func TestAdminCancelEventReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{
		event: &domain.Event{ID: 1, Status: domain.EventStatusCancelled},
	})
	router := gin.New()
	router.DELETE("/admin/eventos/:id", controller.Cancel)

	request := httptest.NewRequest(http.MethodDelete, "/admin/eventos/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAdminReportEventReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{
		report: &services.EventReport{EventID: 1, OccupancyPercentage: 18},
	})
	router := gin.New()
	router.GET("/admin/eventos/:id/reporte", controller.Report)

	request := httptest.NewRequest(http.MethodGet, "/admin/eventos/1/reporte", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestAdminEventInvalidIDReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{})
	router := gin.New()
	router.DELETE("/admin/eventos/:id", controller.Cancel)

	request := httptest.NewRequest(http.MethodDelete, "/admin/eventos/invalid", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestAdminUpdateInvalidDateReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewAdminEventController(fakeAdminEventService{})
	router := gin.New()
	router.PATCH("/admin/eventos/:id", controller.Update)

	request := httptest.NewRequest(http.MethodPatch, "/admin/eventos/1", strings.NewReader(`{"date":"invalid"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestParseRequiredDateAcceptsRFC3339(t *testing.T) {
	if _, err := parseRequiredDate(time.Date(2026, 12, 10, 20, 0, 0, 0, time.UTC).Format(time.RFC3339)); err != nil {
		t.Fatalf("expected RFC3339 date to parse: %v", err)
	}
}
