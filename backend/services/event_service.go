package services

import (
	"backend/domain"
)

type EventService struct {
	eventDAO EventRepository
}

type EventRepository interface {
	GetAll() ([]domain.Event, error)
	GetByID(id uint) (*domain.Event, error)
}

func NewEventService(eventDAO EventRepository) *EventService {
	return &EventService{eventDAO: eventDAO}
}

func (service *EventService) ListEvents() ([]domain.Event, error) {
	return service.eventDAO.GetAll()
}

func (service *EventService) GetEventByID(id uint) (*domain.Event, error) {
	return service.eventDAO.GetByID(id)
}
