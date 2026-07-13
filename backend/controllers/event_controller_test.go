package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type fakeEventRepository struct {
	events          []domain.Event
	event           *domain.Event
	err             error
	receivedFilters domain.EventFilters
}

func (repo *fakeEventRepository) GetAll(filters domain.EventFilters) ([]domain.Event, error) {
	repo.receivedFilters = filters
	return repo.events, repo.err
}

func (repo *fakeEventRepository) GetByID(id uint) (*domain.Event, error) {
	if repo.err != nil {
		return nil, repo.err
	}

	return repo.event, nil
}

func TestEventGetAllReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewEventController(services.NewEventService(&fakeEventRepository{
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

func TestEventGetAllWithSearchFilter(t *testing.T) {
	repo := serveEventListRequest(t, "/eventos?search=rock")

	if repo.receivedFilters.Search != "rock" {
		t.Fatalf("expected search rock, got %q", repo.receivedFilters.Search)
	}
}

func TestEventGetAllWithLocationFilter(t *testing.T) {
	repo := serveEventListRequest(t, "/eventos?location=Cordoba")

	if repo.receivedFilters.Location != "Cordoba" {
		t.Fatalf("expected location Cordoba, got %q", repo.receivedFilters.Location)
	}
}

func TestEventGetAllWithDateRangeFilters(t *testing.T) {
	repo := serveEventListRequest(t, "/eventos?date_from=2026-01-01&date_to=2026-01-31")

	if repo.receivedFilters.DateFrom == nil || repo.receivedFilters.DateTo == nil {
		t.Fatal("expected date range filters")
	}

	expectedFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if !repo.receivedFilters.DateFrom.Equal(expectedFrom) {
		t.Fatalf("expected date_from %s, got %s", expectedFrom, repo.receivedFilters.DateFrom)
	}
}

func TestEventGetAllWithPriceRangeFilters(t *testing.T) {
	repo := serveEventListRequest(t, "/eventos?min_price=100&max_price=500")

	if repo.receivedFilters.MinPrice == nil || *repo.receivedFilters.MinPrice != 100 {
		t.Fatalf("expected min_price 100, got %+v", repo.receivedFilters.MinPrice)
	}

	if repo.receivedFilters.MaxPrice == nil || *repo.receivedFilters.MaxPrice != 500 {
		t.Fatalf("expected max_price 500, got %+v", repo.receivedFilters.MaxPrice)
	}
}

func TestEventGetAllWithFestivalFilter(t *testing.T) {
	repo := serveEventListRequest(t, "/eventos?is_festival=true")

	if repo.receivedFilters.IsFestival == nil || !*repo.receivedFilters.IsFestival {
		t.Fatalf("expected is_festival true, got %+v", repo.receivedFilters.IsFestival)
	}
}

func TestEventGetAllWithAvailableFilter(t *testing.T) {
	repo := serveEventListRequest(t, "/eventos?available=true")

	if !repo.receivedFilters.AvailableOnly {
		t.Fatal("expected available filter true")
	}
}

func TestEventGetAllWithSortFilter(t *testing.T) {
	repo := serveEventListRequest(t, "/eventos?sort=price_desc")

	if repo.receivedFilters.Sort != domain.EventSortPriceDesc {
		t.Fatalf("expected sort price_desc, got %q", repo.receivedFilters.Sort)
	}
}

func TestEventGetAllWithCombinedFilters(t *testing.T) {
	repo := serveEventListRequest(t, "/eventos?search=festival&location=Rosario&available=true&sort=date_desc")

	if repo.receivedFilters.Search != "festival" ||
		repo.receivedFilters.Location != "Rosario" ||
		!repo.receivedFilters.AvailableOnly ||
		repo.receivedFilters.Sort != domain.EventSortDateDesc {
		t.Fatalf("unexpected combined filters: %+v", repo.receivedFilters)
	}
}

func TestEventGetAllWithInvalidParametersReturnsBadRequest(t *testing.T) {
	tests := []string{
		"/eventos?date_from=invalid",
		"/eventos?date_to=invalid",
		"/eventos?min_price=abc",
		"/eventos?max_price=-1",
		"/eventos?is_festival=yes",
		"/eventos?available=1",
		"/eventos?sort=random",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			controller := NewEventController(services.NewEventService(&fakeEventRepository{}))
			router := gin.New()
			router.GET("/eventos", controller.GetAll)

			request := httptest.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", response.Code)
			}
		})
	}
}

func TestEventGetAllDoesNotReturnCancelledEventsFromRepository(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewEventController(services.NewEventService(&fakeEventRepository{
		events: []domain.Event{{ID: 1, Status: domain.EventStatusActive}},
	}))
	router := gin.New()
	router.GET("/eventos", controller.GetAll)

	request := httptest.NewRequest(http.MethodGet, "/eventos", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	var body struct {
		Data []domain.Event `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	for _, event := range body.Data {
		if event.Status == domain.EventStatusCancelled {
			t.Fatalf("expected no cancelled events, got %+v", event)
		}
	}
}

func TestEventGetByIDWithInvalidIDReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := NewEventController(services.NewEventService(&fakeEventRepository{}))
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

	controller := NewEventController(services.NewEventService(&fakeEventRepository{err: gorm.ErrRecordNotFound}))
	router := gin.New()
	router.GET("/eventos/:id", controller.GetByID)

	request := httptest.NewRequest(http.MethodGet, "/eventos/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func serveEventListRequest(t *testing.T, path string) *fakeEventRepository {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := &fakeEventRepository{events: []domain.Event{}}
	controller := NewEventController(services.NewEventService(repo))
	router := gin.New()
	router.GET("/eventos", controller.GetAll)

	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	return repo
}
