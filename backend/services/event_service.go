package services

import (
	"backend/domain"
)

type EventService struct {
	eventDAO EventRepository
}

type EventRepository interface {
	GetAll(filters domain.EventFilters) ([]domain.Event, error)
	GetByID(id uint) (*domain.Event, error)
}

func NewEventService(eventDAO EventRepository) *EventService {
	return &EventService{eventDAO: eventDAO}
}

func (service *EventService) ListEvents(filters domain.EventFilters) ([]domain.Event, error) {
	return service.eventDAO.GetAll(filters)
}

func (service *EventService) GetEventByID(id uint) (*domain.Event, error) {
	return service.eventDAO.GetByID(id)
}
