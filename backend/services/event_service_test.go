package services

import (
	"errors"
	"testing"

	"backend/domain"

	"gorm.io/gorm"
)

type fakeEventRepository struct {
	events          []domain.Event
	event           *domain.Event
	err             error
	receivedFilters domain.EventFilters
}

func (repo *fakeEventRepository) GetAll(filters domain.EventFilters) ([]domain.Event, error) {
	repo.receivedFilters = filters
	return repo.events, repo.err
}

func (repo *fakeEventRepository) GetByID(id uint) (*domain.Event, error) {
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
	repo := &fakeEventRepository{events: expectedEvents}
	service := NewEventService(repo)

	events, err := service.ListEvents(domain.EventFilters{})
	if err != nil {
		t.Fatalf("ListEvents returned error: %v", err)
	}

	if len(events) != len(expectedEvents) {
		t.Fatalf("expected %d events, got %d", len(expectedEvents), len(events))
	}
}

func TestListEventsPassesFilters(t *testing.T) {
	repo := &fakeEventRepository{}
	service := NewEventService(repo)
	isFestival := true
	filters := domain.EventFilters{
		Search:     "rock",
		IsFestival: &isFestival,
		Sort:       domain.EventSortPriceDesc,
	}

	if _, err := service.ListEvents(filters); err != nil {
		t.Fatalf("ListEvents returned error: %v", err)
	}

	if repo.receivedFilters.Search != filters.Search || repo.receivedFilters.Sort != filters.Sort {
		t.Fatalf("expected filters to be passed, got %+v", repo.receivedFilters)
	}

	if repo.receivedFilters.IsFestival == nil || !*repo.receivedFilters.IsFestival {
		t.Fatalf("expected is_festival filter true, got %+v", repo.receivedFilters.IsFestival)
	}
}

func TestGetEventByIDWithMissingEventReturnsError(t *testing.T) {
	service := NewEventService(&fakeEventRepository{err: gorm.ErrRecordNotFound})

	_, err := service.GetEventByID(1)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}
