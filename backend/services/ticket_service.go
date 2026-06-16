package services

import (
	"errors"
	"fmt"
	"time"

	"backend/dao"
	"backend/domain"
)

var ErrEventSoldOut = errors.New("event capacity exhausted")

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

func (service *TicketService) generateTicketCode(now time.Time) (string, error) {
	count, err := service.ticketDAO.CountTickets()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("TICKET-%d-%06d", now.Year(), count+1), nil
}
