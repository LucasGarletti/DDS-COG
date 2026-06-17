package services

import (
	"errors"
	"fmt"
	"time"

	"backend/dao"
	"backend/domain"
)

var (
	ErrEventSoldOut           = errors.New("event capacity exhausted")
	ErrTicketNotOwned         = errors.New("ticket does not belong to authenticated user")
	ErrTicketAlreadyCancelled = errors.New("ticket already cancelled")
)

type TicketService struct {
	ticketDAO *dao.TicketDAO
}

type PurchaseTicketInput struct {
	UserID  uint
	EventID uint
}

func NewTicketService(ticketDAO *dao.TicketDAO) *TicketService {
	return &TicketService{ticketDAO: ticketDAO}
}

func (service *TicketService) PurchaseTicket(input PurchaseTicketInput) (*domain.Ticket, error) {
	event, err := service.ticketDAO.GetEventByID(input.EventID)
	if err != nil {
		return nil, err
	}

	if event.AvailableCapacity <= 0 {
		return nil, ErrEventSoldOut
	}

	now := time.Now()
	code, err := service.generateTicketCode(now)
	if err != nil {
		return nil, err
	}

	ticket := &domain.Ticket{
		UserID:       input.UserID,
		EventID:      input.EventID,
		Code:         code,
		Status:       domain.TicketStatusActive,
		PurchaseDate: now,
	}

	event.AvailableCapacity--

	if err := service.ticketDAO.CreatePurchase(ticket, event); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (service *TicketService) ListTicketsByUser(userID uint) ([]domain.Ticket, error) {
	return service.ticketDAO.GetByUserID(userID)
}

func (service *TicketService) CancelTicket(userID uint, ticketID uint) (*domain.Ticket, error) {
	ticket, err := service.ticketDAO.GetByIDWithEvent(ticketID)
	if err != nil {
		return nil, err
	}

	if ticket.UserID != userID {
		return nil, ErrTicketNotOwned
	}

	if ticket.Status == domain.TicketStatusCancelled {
		return nil, ErrTicketAlreadyCancelled
	}

	now := time.Now()
	ticket.Status = domain.TicketStatusCancelled
	ticket.CancellationDate = &now
	ticket.Event.AvailableCapacity++

	if err := service.ticketDAO.SaveTicketAndEvent(ticket, &ticket.Event); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (service *TicketService) generateTicketCode(now time.Time) (string, error) {
	count, err := service.ticketDAO.CountTickets()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("TICKET-%d-%06d", now.Year(), count+1), nil
}
