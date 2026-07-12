package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type fakeUserItineraryService struct {
	itinerary *services.ItineraryResult
	addResult *services.AddItineraryItemResult
	err       error
}

func (service fakeUserItineraryService) GetItinerary(userID uint, eventID uint) (*services.ItineraryResult, error) {
	return service.itinerary, service.err
}

func (service fakeUserItineraryService) AddShow(userID uint, eventID uint, scheduleID uint) (*services.AddItineraryItemResult, error) {
	return service.addResult, service.err
}

func (service fakeUserItineraryService) AddPersonalActivity(input services.AddPersonalActivityInput) (*services.AddItineraryItemResult, error) {
	return service.addResult, service.err
}

func (service fakeUserItineraryService) DeleteItem(userID uint, eventID uint, itemID uint) error {
	return service.err
}

func TestUserItineraryGetReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewUserItineraryController(fakeUserItineraryService{
		itinerary: &services.ItineraryResult{Event: &domain.Event{ID: 1}, Items: []domain.ItineraryItem{}},
	})
	router := gin.New()
	router.GET("/mis-itinerarios/:eventoId", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Get(c)
	})

	request := httptest.NewRequest(http.MethodGet, "/mis-itinerarios/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestUserItineraryAddShowForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewUserItineraryController(fakeUserItineraryService{err: services.ErrItineraryForbidden})
	router := gin.New()
	router.POST("/mis-itinerarios/:eventoId/shows", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.AddShow(c)
	})

	request := httptest.NewRequest(http.MethodPost, "/mis-itinerarios/1/shows", strings.NewReader(`{"festival_schedule_id":1}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
}

func TestUserItineraryAddPersonalNonFestivalReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewUserItineraryController(fakeUserItineraryService{err: services.ErrEventNotFestival})
	router := gin.New()
	router.POST("/mis-itinerarios/:eventoId/actividades", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.AddPersonalActivity(c)
	})

	body := `{"title":"Comer","start_time":"2026-12-10T20:00:00Z","end_time":"2026-12-10T21:00:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/mis-itinerarios/1/actividades", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestUserItineraryAddPersonalInvalidTimeReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewUserItineraryController(fakeUserItineraryService{})
	router := gin.New()
	router.POST("/mis-itinerarios/:eventoId/actividades", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.AddPersonalActivity(c)
	})

	body := `{"title":"Comer","start_time":"bad","end_time":"2026-12-10T21:00:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/mis-itinerarios/1/actividades", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestUserItineraryDeleteItemReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewUserItineraryController(fakeUserItineraryService{})
	router := gin.New()
	router.DELETE("/mis-itinerarios/:eventoId/items/:itemId", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.DeleteItem(c)
	})

	request := httptest.NewRequest(http.MethodDelete, "/mis-itinerarios/1/items/2", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestUserItineraryInternalErrorReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewUserItineraryController(fakeUserItineraryService{err: errors.New("boom")})
	router := gin.New()
	router.GET("/mis-itinerarios/:eventoId", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Get(c)
	})

	request := httptest.NewRequest(http.MethodGet, "/mis-itinerarios/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
}
