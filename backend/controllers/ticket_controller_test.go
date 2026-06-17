package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type fakeTicketRepository struct {
	tickets []domain.Ticket
}

func (repo fakeTicketRepository) GetEventByID(id uint) (*domain.Event, error) {
	return &domain.Event{ID: id, AvailableCapacity: 10}, nil
}

func (repo fakeTicketRepository) CountTickets() (int64, error) {
	return 0, nil
}

func (repo fakeTicketRepository) GetByUserID(userID uint) ([]domain.Ticket, error) {
	return repo.tickets, nil
}

func (repo fakeTicketRepository) GetByIDWithEvent(id uint) (*domain.Ticket, error) {
	return &domain.Ticket{ID: id, UserID: 1, Status: domain.TicketStatusActive}, nil
}

func (repo fakeTicketRepository) CreatePurchase(ticket *domain.Ticket, event *domain.Event) error {
	return nil
}

func (repo fakeTicketRepository) SaveTicketAndEvent(ticket *domain.Ticket, event *domain.Event) error {
	return nil
}

func (repo fakeTicketRepository) SaveTicket(ticket *domain.Ticket) error {
	return nil
}

type fakeUserRepository struct{}

func (repo fakeUserRepository) FindByEmail(email string) (*domain.User, error) {
	return &domain.User{ID: 2, Email: email}, nil
}

func (repo fakeUserRepository) Create(user *domain.User) error {
	return nil
}

func newTicketTestController(tickets []domain.Ticket) *TicketController {
	return NewTicketController(services.NewTicketService(
		fakeTicketRepository{tickets: tickets},
		fakeUserRepository{},
	))
}

func TestTicketPurchaseWithInvalidEventIDReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestController(nil)
	router := gin.New()
	router.POST("/entradas/comprar/:eventoId", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Purchase(c)
	})

	request := httptest.NewRequest(http.MethodPost, "/entradas/comprar/invalid", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestTicketCancelWithInvalidIDReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestController(nil)
	router := gin.New()
	router.PATCH("/entradas/:id/cancelar", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Cancel(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/invalid/cancelar", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestTicketTransferWithInvalidBodyReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestController(nil)
	router := gin.New()
	router.PATCH("/entradas/:id/transferir", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Transfer(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/1/transferir", strings.NewReader("{"))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestGetMyTicketsWithUserInContextReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestController([]domain.Ticket{{ID: 1, UserID: 1, Status: domain.TicketStatusActive}})
	router := gin.New()
	router.GET("/mis-entradas", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.GetMyTickets(c)
	})

	request := httptest.NewRequest(http.MethodGet, "/mis-entradas", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}
