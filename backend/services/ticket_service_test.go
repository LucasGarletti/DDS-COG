package services

import (
	"errors"
	"testing"

	"backend/domain"
)

type fakeTicketRepository struct {
	event  *domain.Event
	ticket *domain.Ticket
}

func (repo fakeTicketRepository) GetEventByID(id uint) (*domain.Event, error) {
	return repo.event, nil
}

func (repo fakeTicketRepository) CountTickets() (int64, error) {
	return 0, nil
}

func (repo fakeTicketRepository) GetByUserID(userID uint) ([]domain.Ticket, error) {
	return nil, nil
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
