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
	"gorm.io/gorm"
)

type fakeTicketRepository struct {
	event      *domain.Event
	ticket     *domain.Ticket
	tickets    []domain.Ticket
	eventErr   error
	ticketErr  error
	ticketsErr error
}

func (repo fakeTicketRepository) GetEventByID(id uint) (*domain.Event, error) {
	if repo.eventErr != nil {
		return nil, repo.eventErr
	}
	if repo.event != nil {
		return repo.event, nil
	}
	return &domain.Event{ID: id, AvailableCapacity: 10}, nil
}

func (repo fakeTicketRepository) CountTickets() (int64, error) {
	return 0, nil
}

func (repo fakeTicketRepository) GetByUserID(userID uint) ([]domain.Ticket, error) {
	if repo.ticketsErr != nil {
		return nil, repo.ticketsErr
	}
	return repo.tickets, nil
}

func (repo fakeTicketRepository) GetByIDWithEvent(id uint) (*domain.Ticket, error) {
	if repo.ticketErr != nil {
		return nil, repo.ticketErr
	}
	if repo.ticket != nil {
		return repo.ticket, nil
	}
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

func (repo fakeTicketRepository) DeleteUserItineraryByUserAndEvent(userID uint, eventID uint) error {
	return nil
}

type fakeUserRepository struct {
	user *domain.User
	err  error
}

func (repo fakeUserRepository) FindByEmail(email string) (*domain.User, error) {
	if repo.err != nil {
		return nil, repo.err
	}
	if repo.user != nil {
		return repo.user, nil
	}
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

func newTicketTestControllerWithRepos(ticketRepo services.TicketRepository, userRepo services.UserRepository) *TicketController {
	return NewTicketController(services.NewTicketService(ticketRepo, userRepo))
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

func TestGetMyTicketsWithoutUserReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestController(nil)
	router := gin.New()
	router.GET("/mis-entradas", controller.GetMyTickets)

	request := httptest.NewRequest(http.MethodGet, "/mis-entradas", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", response.Code)
	}
}

func TestGetMyTicketsWithRepositoryErrorReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{ticketsErr: errors.New("repo error")},
		fakeUserRepository{},
	)
	router := gin.New()
	router.GET("/mis-entradas", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.GetMyTickets(c)
	})

	request := httptest.NewRequest(http.MethodGet, "/mis-entradas", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
}

func TestTicketPurchaseSuccessReturnsCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{event: &domain.Event{ID: 1, AvailableCapacity: 2}},
		fakeUserRepository{},
	)
	router := gin.New()
	router.POST("/entradas/comprar/:eventoId", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Purchase(c)
	})

	request := httptest.NewRequest(http.MethodPost, "/entradas/comprar/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}
}

func TestTicketPurchaseMissingEventReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{eventErr: gorm.ErrRecordNotFound},
		fakeUserRepository{},
	)
	router := gin.New()
	router.POST("/entradas/comprar/:eventoId", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Purchase(c)
	})

	request := httptest.NewRequest(http.MethodPost, "/entradas/comprar/99", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestTicketPurchaseSoldOutReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{event: &domain.Event{ID: 1, AvailableCapacity: 0}},
		fakeUserRepository{},
	)
	router := gin.New()
	router.POST("/entradas/comprar/:eventoId", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Purchase(c)
	})

	request := httptest.NewRequest(http.MethodPost, "/entradas/comprar/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestTicketCancelSuccessReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{ticket: &domain.Ticket{ID: 1, UserID: 1, Status: domain.TicketStatusActive}},
		fakeUserRepository{},
	)
	router := gin.New()
	router.PATCH("/entradas/:id/cancelar", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Cancel(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/1/cancelar", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestTicketCancelMissingTicketReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{ticketErr: gorm.ErrRecordNotFound},
		fakeUserRepository{},
	)
	router := gin.New()
	router.PATCH("/entradas/:id/cancelar", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Cancel(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/99/cancelar", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestTicketCancelAlreadyCancelledReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{ticket: &domain.Ticket{ID: 1, UserID: 1, Status: domain.TicketStatusCancelled}},
		fakeUserRepository{},
	)
	router := gin.New()
	router.PATCH("/entradas/:id/cancelar", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Cancel(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/1/cancelar", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestTicketTransferSuccessReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{ticket: &domain.Ticket{ID: 1, UserID: 1, Status: domain.TicketStatusActive}},
		fakeUserRepository{user: &domain.User{ID: 2, Email: "other@mail.com"}},
	)
	router := gin.New()
	router.PATCH("/entradas/:id/transferir", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Transfer(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/1/transferir", strings.NewReader(`{"email_destinatario":"other@mail.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestTicketTransferMissingRecipientReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{ticket: &domain.Ticket{ID: 1, UserID: 1, Status: domain.TicketStatusActive}},
		fakeUserRepository{err: gorm.ErrRecordNotFound},
	)
	router := gin.New()
	router.PATCH("/entradas/:id/transferir", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Transfer(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/1/transferir", strings.NewReader(`{"email_destinatario":"missing@mail.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestTicketTransferToSameUserReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{ticket: &domain.Ticket{ID: 1, UserID: 1, Status: domain.TicketStatusActive}},
		fakeUserRepository{user: &domain.User{ID: 1, Email: "same@mail.com"}},
	)
	router := gin.New()
	router.PATCH("/entradas/:id/transferir", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Transfer(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/1/transferir", strings.NewReader(`{"email_destinatario":"same@mail.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestTicketTransferCancelledTicketReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := newTicketTestControllerWithRepos(
		fakeTicketRepository{ticket: &domain.Ticket{ID: 1, UserID: 1, Status: domain.TicketStatusCancelled}},
		fakeUserRepository{},
	)
	router := gin.New()
	router.PATCH("/entradas/:id/transferir", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Transfer(c)
	})

	request := httptest.NewRequest(http.MethodPatch, "/entradas/1/transferir", strings.NewReader(`{"email_destinatario":"other@mail.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}
