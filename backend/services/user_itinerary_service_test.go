package services

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

type fakeUserItineraryRepository struct {
	hasActiveTicket bool
	itinerary       *domain.UserItinerary
	items           []domain.ItineraryItem
	item            *domain.ItineraryItem
	err             error
	createdItem     *domain.ItineraryItem
	deletedItemID   uint
}

func (repo *fakeUserItineraryRepository) HasActiveTicket(userID uint, eventID uint) (bool, error) {
	return repo.hasActiveTicket, repo.err
}

func (repo *fakeUserItineraryRepository) GetByUserAndEvent(userID uint, eventID uint) (*domain.UserItinerary, error) {
	if repo.itinerary == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return repo.itinerary, repo.err
}

func (repo *fakeUserItineraryRepository) Create(itinerary *domain.UserItinerary) error {
	itinerary.ID = 1
	repo.itinerary = itinerary
	return repo.err
}

func (repo *fakeUserItineraryRepository) GetItemsByItineraryID(itineraryID uint) ([]domain.ItineraryItem, error) {
	return repo.items, repo.err
}

func (repo *fakeUserItineraryRepository) CreateItem(item *domain.ItineraryItem) error {
	repo.createdItem = item
	return repo.err
}

func (repo *fakeUserItineraryRepository) GetItemByID(id uint) (*domain.ItineraryItem, error) {
	if repo.err != nil {
		return nil, repo.err
	}
	return repo.item, nil
}

func (repo *fakeUserItineraryRepository) DeleteItem(id uint) error {
	repo.deletedItemID = id
	return repo.err
}

func newUserItineraryTestService(itineraryRepo *fakeUserItineraryRepository, scheduleRepo *fakeFestivalScheduleRepository, isFestival bool) *UserItineraryService {
	return NewUserItineraryService(
		fakeFestivalEventRepository{event: &domain.Event{ID: 1, IsFestival: isFestival}},
		scheduleRepo,
		itineraryRepo,
	)
}

func TestUserItineraryAddShowCopiesScheduleData(t *testing.T) {
	start := time.Date(2026, 12, 10, 20, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	itineraryRepo := &fakeUserItineraryRepository{hasActiveTicket: true, itinerary: &domain.UserItinerary{ID: 1, UserID: 1, EventID: 1}}
	scheduleRepo := &fakeFestivalScheduleRepository{schedule: &domain.FestivalSchedule{ID: 3, EventID: 1, Artist: "Banda", Stage: "Norte", StartTime: start, EndTime: end}}
	service := newUserItineraryTestService(itineraryRepo, scheduleRepo, true)

	result, err := service.AddShow(1, 1, 3)
	if err != nil {
		t.Fatalf("AddShow returned error: %v", err)
	}

	if result.Item.Title != "Banda" || result.Item.StartTime != start || result.Item.EndTime != end {
		t.Fatalf("expected show data copied from schedule, got %+v", result.Item)
	}
}

func TestUserItineraryAddPersonalDetectsOverlap(t *testing.T) {
	start := time.Date(2026, 12, 10, 20, 0, 0, 0, time.UTC)
	itineraryRepo := &fakeUserItineraryRepository{
		hasActiveTicket: true,
		itinerary:       &domain.UserItinerary{ID: 1, UserID: 1, EventID: 1},
		items:           []domain.ItineraryItem{{StartTime: start, EndTime: start.Add(time.Hour)}},
	}
	service := newUserItineraryTestService(itineraryRepo, &fakeFestivalScheduleRepository{}, true)

	result, err := service.AddPersonalActivity(AddPersonalActivityInput{
		UserID:    1,
		EventID:   1,
		Title:     "Comer",
		StartTime: start.Add(30 * time.Minute),
		EndTime:   start.Add(90 * time.Minute),
	})
	if err != nil {
		t.Fatalf("AddPersonalActivity returned error: %v", err)
	}

	if !result.HasConflict {
		t.Fatal("expected conflict to be detected")
	}
}

func TestUserItineraryNonFestivalReturnsError(t *testing.T) {
	service := newUserItineraryTestService(&fakeUserItineraryRepository{hasActiveTicket: true}, &fakeFestivalScheduleRepository{}, false)

	_, err := service.GetItinerary(1, 1)
	if !errors.Is(err, ErrEventNotFestival) {
		t.Fatalf("expected ErrEventNotFestival, got %v", err)
	}
}

func TestUserItineraryWithoutActiveTicketReturnsForbidden(t *testing.T) {
	service := newUserItineraryTestService(&fakeUserItineraryRepository{hasActiveTicket: false}, &fakeFestivalScheduleRepository{}, true)

	_, err := service.GetItinerary(1, 1)
	if !errors.Is(err, ErrItineraryForbidden) {
		t.Fatalf("expected ErrItineraryForbidden, got %v", err)
	}
}

func TestUserItineraryCancelledTicketReturnsForbidden(t *testing.T) {
	service := newUserItineraryTestService(&fakeUserItineraryRepository{hasActiveTicket: false}, &fakeFestivalScheduleRepository{}, true)

	_, err := service.AddPersonalActivity(AddPersonalActivityInput{
		UserID:    1,
		EventID:   1,
		Title:     "Comer",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})
	if !errors.Is(err, ErrItineraryForbidden) {
		t.Fatalf("expected ErrItineraryForbidden, got %v", err)
	}
}

func TestUserItineraryGetReportsConflicts(t *testing.T) {
	start := time.Date(2026, 12, 10, 20, 0, 0, 0, time.UTC)
	service := newUserItineraryTestService(&fakeUserItineraryRepository{
		hasActiveTicket: true,
		itinerary:       &domain.UserItinerary{ID: 1, UserID: 1, EventID: 1},
		items: []domain.ItineraryItem{
			{StartTime: start, EndTime: start.Add(time.Hour)},
			{StartTime: start.Add(30 * time.Minute), EndTime: start.Add(90 * time.Minute)},
		},
	}, &fakeFestivalScheduleRepository{}, true)

	result, err := service.GetItinerary(1, 1)
	if err != nil {
		t.Fatalf("GetItinerary returned error: %v", err)
	}

	if !result.HasConflicts {
		t.Fatal("expected has_conflicts true")
	}
}

func TestUserItineraryDeleteItemNotOwned(t *testing.T) {
	service := newUserItineraryTestService(&fakeUserItineraryRepository{
		hasActiveTicket: true,
		itinerary:       &domain.UserItinerary{ID: 1, UserID: 1, EventID: 1},
		item:            &domain.ItineraryItem{ID: 9, UserItineraryID: 2},
	}, &fakeFestivalScheduleRepository{}, true)

	err := service.DeleteItem(1, 1, 9)
	if !errors.Is(err, ErrItineraryItemNotOwned) {
		t.Fatalf("expected ErrItineraryItemNotOwned, got %v", err)
	}
}
