package services

import (
	"errors"
	"testing"

	"backend/domain"

	"gorm.io/gorm"
)

type fakeEventRepository struct {
	events []domain.Event
	event  *domain.Event
	err    error
}

func (repo fakeEventRepository) GetAll() ([]domain.Event, error) {
	return repo.events, repo.err
}

func (repo fakeEventRepository) GetByID(id uint) (*domain.Event, error) {
	if repo.err != nil {
		return nil, repo.err
	}

	return repo.event, nil
}

func TestListEventsReturnsList(t *testing.T) {
	expectedEvents := []domain.Event{
		{ID: 1, Title: "Event 1"},
		{ID: 2, Title: "Event 2"},
	}
	service := NewEventService(fakeEventRepository{events: expectedEvents})

	events, err := service.ListEvents()
	if err != nil {
		t.Fatalf("ListEvents returned error: %v", err)
	}

	if len(events) != len(expectedEvents) {
		t.Fatalf("expected %d events, got %d", len(expectedEvents), len(events))
	}
}

func TestGetEventByIDWithMissingEventReturnsError(t *testing.T) {
	service := NewEventService(fakeEventRepository{err: gorm.ErrRecordNotFound})

	_, err := service.GetEventByID(1)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}
