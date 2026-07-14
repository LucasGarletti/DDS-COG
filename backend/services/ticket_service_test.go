package services

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

type fakeTicketRepository struct {
	event                   *domain.Event
	ticket                  *domain.Ticket
	tickets                 []domain.Ticket
	count                   int64
	eventErr                error
	countErr                error
	ticketsErr              error
	ticketErr               error
	createErr               error
	saveTicketAndEventErr   error
	saveTicketErr           error
	deleteItineraryErr      error
	deletedItineraryUserID  uint
	deletedItineraryEventID uint
}

func (repo fakeTicketRepository) GetEventByID(id uint) (*domain.Event, error) {
	if repo.eventErr != nil {
		return nil, repo.eventErr
	}
	return repo.event, nil
}

func (repo fakeTicketRepository) CountTickets() (int64, error) {
	if repo.countErr != nil {
		return 0, repo.countErr
	}
	return repo.count, nil
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
	return repo.ticket, nil
}

func (repo fakeTicketRepository) CreatePurchase(ticket *domain.Ticket, event *domain.Event) error {
	return repo.createErr
}

func (repo fakeTicketRepository) SaveTicketAndEvent(ticket *domain.Ticket, event *domain.Event) error {
	return repo.saveTicketAndEventErr
}

func (repo fakeTicketRepository) SaveTicket(ticket *domain.Ticket) error {
	return repo.saveTicketErr
}

func (repo fakeTicketRepository) DeleteUserItineraryByUserAndEvent(userID uint, eventID uint) error {
	return repo.deleteItineraryErr
}

func TestPurchaseTicketWithSoldOutEventReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		event: &domain.Event{ID: 1, AvailableCapacity: 0},
	}, fakeUserRepository{})

	_, err := service.PurchaseTicket(PurchaseTicketInput{
		UserID:  1,
		EventID: 1,
	})
	if !errors.Is(err, ErrEventSoldOut) {
		t.Fatalf("expected ErrEventSoldOut, got %v", err)
	}
}

func TestPurchaseTicketWithCancelledEventReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		event: &domain.Event{ID: 1, Status: domain.EventStatusCancelled, AvailableCapacity: 10},
	}, fakeUserRepository{})

	_, err := service.PurchaseTicket(PurchaseTicketInput{
		UserID:  1,
		EventID: 1,
	})
	if !errors.Is(err, ErrEventCancelledPurchase) {
		t.Fatalf("expected ErrEventCancelledPurchase, got %v", err)
	}
}

func TestCancelAlreadyCancelledTicketReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{
			ID:     1,
			UserID: 1,
			Status: domain.TicketStatusCancelled,
		},
	}, fakeUserRepository{})

	_, err := service.CancelTicket(1, 1)
	if !errors.Is(err, ErrTicketAlreadyCancelled) {
		t.Fatalf("expected ErrTicketAlreadyCancelled, got %v", err)
	}
}

func TestTransferCancelledTicketReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{
			ID:     1,
			UserID: 1,
			Status: domain.TicketStatusCancelled,
		},
	}, fakeUserRepository{})

	_, err := service.TransferTicket(TransferTicketInput{
		UserID:         1,
		TicketID:       1,
		RecipientEmail: "other@mail.com",
	})
	if !errors.Is(err, ErrTicketAlreadyCancelled) {
		t.Fatalf("expected ErrTicketAlreadyCancelled, got %v", err)
	}
}

func TestTransferToSameUserReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{
			ID:     1,
			UserID: 1,
			Status: domain.TicketStatusActive,
		},
	}, fakeUserRepository{
		user: &domain.User{
			ID:    1,
			Email: "same@mail.com",
		},
	})

	_, err := service.TransferTicket(TransferTicketInput{
		UserID:         1,
		TicketID:       1,
		RecipientEmail: "same@mail.com",
	})
	if !errors.Is(err, ErrInvalidRecipient) {
		t.Fatalf("expected ErrInvalidRecipient, got %v", err)
	}
}

type trackingTicketRepository struct {
	fakeTicketRepository
	deletedUserID  uint
	deletedEventID uint
	savedTicket    *domain.Ticket
}

func (repo *trackingTicketRepository) DeleteUserItineraryByUserAndEvent(userID uint, eventID uint) error {
	repo.deletedUserID = userID
	repo.deletedEventID = eventID
	return nil
}

func (repo *trackingTicketRepository) SaveTicket(ticket *domain.Ticket) error {
	repo.savedTicket = ticket
	return nil
}

func TestTransferTicketDeletesOriginalUserItinerary(t *testing.T) {
	repo := &trackingTicketRepository{
		fakeTicketRepository: fakeTicketRepository{
			ticket: &domain.Ticket{
				ID:      1,
				UserID:  1,
				EventID: 7,
				Status:  domain.TicketStatusActive,
			},
		},
	}
	service := NewTicketService(repo, fakeUserRepository{
		user: &domain.User{ID: 2, Email: "other@mail.com"},
	})

	ticket, err := service.TransferTicket(TransferTicketInput{
		UserID:         1,
		TicketID:       1,
		RecipientEmail: "other@mail.com",
	})
	if err != nil {
		t.Fatalf("TransferTicket returned error: %v", err)
	}

	if repo.deletedUserID != 1 || repo.deletedEventID != 7 {
		t.Fatalf("expected deleted itinerary for user 1 event 7, got user %d event %d", repo.deletedUserID, repo.deletedEventID)
	}

	if ticket.UserID != 2 {
		t.Fatalf("expected transferred ticket user 2, got %d", ticket.UserID)
	}
}

func TestListTicketsByUserReturnsList(t *testing.T) {
	expectedTickets := []domain.Ticket{
		{ID: 1, UserID: 1, Status: domain.TicketStatusActive},
		{ID: 2, UserID: 1, Status: domain.TicketStatusCancelled},
	}
	service := NewTicketService(fakeTicketRepository{tickets: expectedTickets}, fakeUserRepository{})

	tickets, err := service.ListTicketsByUser(1)
	if err != nil {
		t.Fatalf("ListTicketsByUser returned error: %v", err)
	}

	if len(tickets) != len(expectedTickets) {
		t.Fatalf("expected %d tickets, got %d", len(expectedTickets), len(tickets))
	}
}

func TestListTicketsByUserReturnsRepositoryError(t *testing.T) {
	repoErr := errors.New("repo error")
	service := NewTicketService(fakeTicketRepository{ticketsErr: repoErr}, fakeUserRepository{})

	_, err := service.ListTicketsByUser(1)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestPurchaseTicketSuccessGeneratesCode(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		event: &domain.Event{
			ID:                1,
			AvailableCapacity: 10,
		},
		count: 4,
	}, fakeUserRepository{})

	ticket, err := service.PurchaseTicket(PurchaseTicketInput{
		UserID:  1,
		EventID: 1,
	})
	if err != nil {
		t.Fatalf("PurchaseTicket returned error: %v", err)
	}

	if ticket.Code == "" {
		t.Fatal("expected generated ticket code")
	}

	if ticket.Status != domain.TicketStatusActive {
		t.Fatalf("expected active ticket, got %s", ticket.Status)
	}
}

func TestPurchaseTicketWithMissingEventReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{eventErr: gorm.ErrRecordNotFound}, fakeUserRepository{})

	_, err := service.PurchaseTicket(PurchaseTicketInput{UserID: 1, EventID: 99})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func TestPurchaseTicketWithCreateErrorReturnsError(t *testing.T) {
	createErr := errors.New("create error")
	service := NewTicketService(fakeTicketRepository{
		event:     &domain.Event{ID: 1, AvailableCapacity: 1},
		createErr: createErr,
	}, fakeUserRepository{})

	_, err := service.PurchaseTicket(PurchaseTicketInput{UserID: 1, EventID: 1})
	if !errors.Is(err, createErr) {
		t.Fatalf("expected create error, got %v", err)
	}
}

func TestPurchaseTicketWithCountErrorReturnsError(t *testing.T) {
	countErr := errors.New("count error")
	service := NewTicketService(fakeTicketRepository{
		event:    &domain.Event{ID: 1, AvailableCapacity: 1},
		countErr: countErr,
	}, fakeUserRepository{})

	_, err := service.PurchaseTicket(PurchaseTicketInput{UserID: 1, EventID: 1})
	if !errors.Is(err, countErr) {
		t.Fatalf("expected count error, got %v", err)
	}
}

func TestGenerateTicketCodeUsesCurrentYearAndCount(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{count: 41}, fakeUserRepository{})

	code, err := service.generateTicketCode(time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("generateTicketCode returned error: %v", err)
	}

	if code != "TICKET-2026-000042" {
		t.Fatalf("expected generated code TICKET-2026-000042, got %s", code)
	}
}

func TestGenerateTicketCodeReturnsCountError(t *testing.T) {
	countErr := errors.New("count error")
	service := NewTicketService(fakeTicketRepository{countErr: countErr}, fakeUserRepository{})

	_, err := service.generateTicketCode(time.Now())
	if !errors.Is(err, countErr) {
		t.Fatalf("expected count error, got %v", err)
	}
}

func TestCancelTicketSuccessReturnsCapacity(t *testing.T) {
	event := domain.Event{ID: 7, AvailableCapacity: 3}
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{
			ID:      1,
			UserID:  1,
			EventID: 7,
			Status:  domain.TicketStatusActive,
			Event:   event,
		},
	}, fakeUserRepository{})

	ticket, err := service.CancelTicket(1, 1)
	if err != nil {
		t.Fatalf("CancelTicket returned error: %v", err)
	}

	if ticket.Status != domain.TicketStatusCancelled {
		t.Fatalf("expected cancelled ticket, got %s", ticket.Status)
	}
	if ticket.CancellationDate == nil {
		t.Fatal("expected cancellation date")
	}
	if ticket.Event.AvailableCapacity != 4 {
		t.Fatalf("expected capacity 4, got %d", ticket.Event.AvailableCapacity)
	}
}

func TestCancelTicketWithMissingTicketReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{ticketErr: gorm.ErrRecordNotFound}, fakeUserRepository{})

	_, err := service.CancelTicket(1, 99)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func TestCancelTicketNotOwnedReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{ID: 1, UserID: 2, Status: domain.TicketStatusActive},
	}, fakeUserRepository{})

	_, err := service.CancelTicket(1, 1)
	if !errors.Is(err, ErrTicketNotOwned) {
		t.Fatalf("expected ErrTicketNotOwned, got %v", err)
	}
}

func TestCancelTicketWithSaveErrorReturnsError(t *testing.T) {
	saveErr := errors.New("save error")
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{
			ID:     1,
			UserID: 1,
			Status: domain.TicketStatusActive,
			Event:  domain.Event{ID: 7, AvailableCapacity: 3},
		},
		saveTicketAndEventErr: saveErr,
	}, fakeUserRepository{})

	_, err := service.CancelTicket(1, 1)
	if !errors.Is(err, saveErr) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestTransferTicketWithMissingRecipientReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{ID: 1, UserID: 1, Status: domain.TicketStatusActive},
	}, fakeUserRepository{err: gorm.ErrRecordNotFound})

	_, err := service.TransferTicket(TransferTicketInput{
		UserID:         1,
		TicketID:       1,
		RecipientEmail: "missing@mail.com",
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func TestTransferTicketNotOwnedReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{ID: 1, UserID: 2, Status: domain.TicketStatusActive},
	}, fakeUserRepository{})

	_, err := service.TransferTicket(TransferTicketInput{
		UserID:         1,
		TicketID:       1,
		RecipientEmail: "other@mail.com",
	})
	if !errors.Is(err, ErrTicketNotOwned) {
		t.Fatalf("expected ErrTicketNotOwned, got %v", err)
	}
}

func TestTransferTicketWithMissingTicketReturnsError(t *testing.T) {
	service := NewTicketService(fakeTicketRepository{ticketErr: gorm.ErrRecordNotFound}, fakeUserRepository{})

	_, err := service.TransferTicket(TransferTicketInput{
		UserID:         1,
		TicketID:       99,
		RecipientEmail: "other@mail.com",
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func TestTransferTicketWithDeleteItineraryErrorReturnsError(t *testing.T) {
	deleteErr := errors.New("delete error")
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{
			ID:      1,
			UserID:  1,
			EventID: 7,
			Status:  domain.TicketStatusActive,
		},
		deleteItineraryErr: deleteErr,
	}, fakeUserRepository{user: &domain.User{ID: 2, Email: "other@mail.com"}})

	_, err := service.TransferTicket(TransferTicketInput{
		UserID:         1,
		TicketID:       1,
		RecipientEmail: "other@mail.com",
	})
	if !errors.Is(err, deleteErr) {
		t.Fatalf("expected delete error, got %v", err)
	}
}

func TestTransferTicketWithSaveErrorReturnsError(t *testing.T) {
	saveErr := errors.New("save error")
	service := NewTicketService(fakeTicketRepository{
		ticket: &domain.Ticket{
			ID:      1,
			UserID:  1,
			EventID: 7,
			Status:  domain.TicketStatusActive,
		},
		saveTicketErr: saveErr,
	}, fakeUserRepository{user: &domain.User{ID: 2, Email: "other@mail.com"}})

	_, err := service.TransferTicket(TransferTicketInput{
		UserID:         1,
		TicketID:       1,
		RecipientEmail: "other@mail.com",
	})
	if !errors.Is(err, saveErr) {
		t.Fatalf("expected save error, got %v", err)
	}
}
