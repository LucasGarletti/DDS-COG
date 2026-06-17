package services

import (
	"errors"
	"testing"

	"backend/domain"
)

type fakeTicketRepository struct {
	event   *domain.Event
	ticket  *domain.Ticket
	tickets []domain.Ticket
	count   int64
}

func (repo fakeTicketRepository) GetEventByID(id uint) (*domain.Event, error) {
	return repo.event, nil
}

func (repo fakeTicketRepository) CountTickets() (int64, error) {
	return repo.count, nil
}

func (repo fakeTicketRepository) GetByUserID(userID uint) ([]domain.Ticket, error) {
	return repo.tickets, nil
}

func (repo fakeTicketRepository) GetByIDWithEvent(id uint) (*domain.Ticket, error) {
	return repo.ticket, nil
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
