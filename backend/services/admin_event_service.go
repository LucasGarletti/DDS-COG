package services

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"backend/domain"
)

var (
	ErrInvalidEventData      = errors.New("invalid event data")
	ErrEventCancelled        = errors.New("event is cancelled")
	ErrEventCapacityConflict = errors.New("event capacity cannot be lower than issued tickets")
)

type AdminEventRepository interface {
	Create(event *domain.Event) error
	GetByIDForAdmin(id uint) (*domain.Event, error)
	Save(event *domain.Event) error
}

type EventReportRepository interface {
	CountTicketsByEvent(eventID uint) (int64, error)
	CountTicketsByEventAndStatus(eventID uint, status string) (int64, error)
}

type AdminEventService struct {
	eventDAO  AdminEventRepository
	ticketDAO EventReportRepository
	now       func() time.Time
}

type CreateEventInput struct {
	Title       string
	Description string
	Date        time.Time
	Location    string
	Capacity    int
	Price       float64
	ImageURL    string
	IsFestival  bool
}

type UpdateEventInput struct {
	ID          uint
	Title       *string
	Description *string
	Date        *time.Time
	Location    *string
	Capacity    *int
	Price       *float64
	ImageURL    *string
	IsFestival  *bool
}

type EventReport struct {
	EventID             uint    `json:"event_id"`
	Title               string  `json:"title"`
	Status              string  `json:"status"`
	Capacity            int     `json:"capacity"`
	AvailableCapacity   int     `json:"available_capacity"`
	TicketsIssued       int64   `json:"tickets_issued"`
	ActiveTickets       int64   `json:"active_tickets"`
	CancelledTickets    int64   `json:"cancelled_tickets"`
	OccupancyPercentage float64 `json:"occupancy_percentage"`
	EstimatedRevenue    float64 `json:"estimated_revenue"`
}

func NewAdminEventService(eventDAO AdminEventRepository, ticketDAO EventReportRepository) *AdminEventService {
	return &AdminEventService{
		eventDAO:  eventDAO,
		ticketDAO: ticketDAO,
		now:       time.Now,
	}
}

func (service *AdminEventService) CreateEvent(input CreateEventInput) (*domain.Event, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Location = strings.TrimSpace(input.Location)
	input.ImageURL = strings.TrimSpace(input.ImageURL)

	if err := service.validateCreateInput(input); err != nil {
		return nil, err
	}

	event := &domain.Event{
		Title:             input.Title,
		Description:       input.Description,
		Date:              input.Date,
		Location:          input.Location,
		Capacity:          input.Capacity,
		AvailableCapacity: input.Capacity,
		Price:             input.Price,
		ImageURL:          input.ImageURL,
		IsFestival:        input.IsFestival,
		Status:            domain.EventStatusActive,
	}

	if err := service.eventDAO.Create(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (service *AdminEventService) UpdateEvent(input UpdateEventInput) (*domain.Event, error) {
	if input.ID == 0 {
		return nil, ErrInvalidEventData
	}

	event, err := service.eventDAO.GetByIDForAdmin(input.ID)
	if err != nil {
		return nil, err
	}

	if event.Status == domain.EventStatusCancelled {
		return nil, ErrEventCancelled
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, ErrInvalidEventData
		}
		event.Title = title
	}

	if input.Description != nil {
		event.Description = *input.Description
	}

	if input.Date != nil {
		if input.Date.IsZero() || !input.Date.After(service.now()) {
			return nil, ErrInvalidEventData
		}
		event.Date = *input.Date
	}

	if input.Location != nil {
		location := strings.TrimSpace(*input.Location)
		if location == "" {
			return nil, ErrInvalidEventData
		}
		event.Location = location
	}

	if input.Price != nil {
		if *input.Price < 0 {
			return nil, ErrInvalidEventData
		}
		event.Price = *input.Price
	}

	if input.ImageURL != nil {
		imageURL := strings.TrimSpace(*input.ImageURL)
		if !isValidOptionalHTTPURL(imageURL) {
			return nil, ErrInvalidEventData
		}
		event.ImageURL = imageURL
	}

	if input.Capacity != nil {
		if *input.Capacity <= 0 {
			return nil, ErrInvalidEventData
		}

		ticketsIssued := event.Capacity - event.AvailableCapacity
		if ticketsIssued < 0 {
			ticketsIssued = 0
		}

		if *input.Capacity < ticketsIssued {
			return nil, ErrEventCapacityConflict
		}

		event.Capacity = *input.Capacity
		event.AvailableCapacity = *input.Capacity - ticketsIssued
	}

	if input.IsFestival != nil {
		event.IsFestival = *input.IsFestival
	}

	if err := service.eventDAO.Save(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (service *AdminEventService) CancelEvent(id uint) (*domain.Event, error) {
	if id == 0 {
		return nil, ErrInvalidEventData
	}

	event, err := service.eventDAO.GetByIDForAdmin(id)
	if err != nil {
		return nil, err
	}

	if event.Status == domain.EventStatusCancelled {
		return nil, ErrEventCancelled
	}

	event.Status = domain.EventStatusCancelled

	if err := service.eventDAO.Save(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (service *AdminEventService) GetEventReport(id uint) (*EventReport, error) {
	if id == 0 {
		return nil, ErrInvalidEventData
	}

	event, err := service.eventDAO.GetByIDForAdmin(id)
	if err != nil {
		return nil, err
	}

	ticketsIssued, err := service.ticketDAO.CountTicketsByEvent(id)
	if err != nil {
		return nil, err
	}

	activeTickets, err := service.ticketDAO.CountTicketsByEventAndStatus(id, domain.TicketStatusActive)
	if err != nil {
		return nil, err
	}

	cancelledTickets, err := service.ticketDAO.CountTicketsByEventAndStatus(id, domain.TicketStatusCancelled)
	if err != nil {
		return nil, err
	}

	occupancyPercentage := 0.0
	if event.Capacity > 0 {
		occupancyPercentage = float64(activeTickets) / float64(event.Capacity) * 100
	}

	return &EventReport{
		EventID:             event.ID,
		Title:               event.Title,
		Status:              event.Status,
		Capacity:            event.Capacity,
		AvailableCapacity:   event.AvailableCapacity,
		TicketsIssued:       ticketsIssued,
		ActiveTickets:       activeTickets,
		CancelledTickets:    cancelledTickets,
		OccupancyPercentage: occupancyPercentage,
		EstimatedRevenue:    float64(activeTickets) * event.Price,
	}, nil
}

func (service *AdminEventService) validateCreateInput(input CreateEventInput) error {
	if input.Title == "" ||
		input.Location == "" ||
		input.Date.IsZero() ||
		!input.Date.After(service.now()) ||
		input.Capacity <= 0 ||
		input.Price < 0 ||
		!isValidOptionalHTTPURL(input.ImageURL) {
		return ErrInvalidEventData
	}

	return nil
}

func isValidOptionalHTTPURL(value string) bool {
	if value == "" {
		return true
	}

	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}

	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}
