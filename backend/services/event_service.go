package services

import (
	"backend/dao"
	"backend/domain"
)

type EventService struct {
	eventDAO *dao.EventDAO
}

func NewEventService(eventDAO *dao.EventDAO) *EventService {
	return &EventService{eventDAO: eventDAO}
}

func (service *EventService) ListEvents() ([]domain.Event, error) {
	return service.eventDAO.GetAll()
}

func (service *EventService) GetEventByID(id uint) (*domain.Event, error) {
	return service.eventDAO.GetByID(id)
}
