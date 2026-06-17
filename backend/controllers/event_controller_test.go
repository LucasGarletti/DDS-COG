package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type fakeEventRepository struct {
	events []domain.Event
	event  *domain.Event
	err    error
}

func (repo fakeEventRepository) GetAll() ([]domain.Event, error) {
	return repo.events, repo.err
}

func (repo fakeEventRepository) GetByID(id uint) (*domain.Event, error) {
	if repo.err != nil {
		return nil, repo.err
	}

	return repo.event, nil
}

func TestEventGetAllReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewEventController(services.NewEventService(fakeEventRepository{
		events: []domain.Event{{ID: 1, Title: "Event 1"}},
	}))
	router := gin.New()
	router.GET("/eventos", controller.GetAll)

	request := httptest.NewRequest(http.MethodGet, "/eventos", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestEventGetByIDWithInvalidIDReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewEventController(services.NewEventService(fakeEventRepository{}))
	router := gin.New()
	router.GET("/eventos/:id", controller.GetByID)

	request := httptest.NewRequest(http.MethodGet, "/eventos/invalid", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestEventGetByIDWithMissingEventReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewEventController(services.NewEventService(fakeEventRepository{err: gorm.ErrRecordNotFound}))
	router := gin.New()
	router.GET("/eventos/:id", controller.GetByID)

	request := httptest.NewRequest(http.MethodGet, "/eventos/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}
