package dao

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

func TestNewEventDAO(t *testing.T) {
	db := newDAOTestDB(t)

	eventDAO := NewEventDAO(db)

	if eventDAO == nil {
		t.Fatal("expected event DAO")
	}
}

func TestEventDAOCreateSaveAndGetByID(t *testing.T) {
	db := newDAOTestDB(t)
	eventDAO := NewEventDAO(db)
	event := &domain.Event{
		Title:             "Rock",
		Description:       "Rock event",
		Date:              time.Date(2027, 1, 10, 20, 0, 0, 0, time.UTC),
		Location:          "Cordoba",
		Capacity:          100,
		AvailableCapacity: 100,
		Price:             500,
		Status:            domain.EventStatusActive,
	}

	if err := eventDAO.Create(event); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	event.Title = "Rock Updated"
	if err := eventDAO.Save(event); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	found, err := eventDAO.GetByID(event.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if found.Title != "Rock Updated" {
		t.Fatalf("expected updated title, got %s", found.Title)
	}
}

func TestEventDAOGetByIDHidesCancelledEvents(t *testing.T) {
	db := newDAOTestDB(t)
	eventDAO := NewEventDAO(db)
	event := createDAOTestEvent(t, db, "Cancelled", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	event.Status = domain.EventStatusCancelled
	if err := db.Save(&event).Error; err != nil {
		t.Fatalf("could not cancel event: %v", err)
	}

	_, err := eventDAO.GetByID(event.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}

	adminEvent, err := eventDAO.GetByIDForAdmin(event.ID)
	if err != nil {
		t.Fatalf("GetByIDForAdmin returned error: %v", err)
	}
	if adminEvent.Status != domain.EventStatusCancelled {
		t.Fatalf("expected cancelled event for admin, got %+v", adminEvent)
	}
}

func TestEventDAOGetAllAppliesFilters(t *testing.T) {
	db := newDAOTestDB(t)
	eventDAO := NewEventDAO(db)
	rock := createDAOTestEvent(t, db, "Rock Night", time.Date(2027, 1, 10, 20, 0, 0, 0, time.UTC), 100, true)
	rock.Description = "Festival de rock"
	rock.Location = "Cordoba"
	rock.Price = 500
	rock.AvailableCapacity = 10
	if err := db.Save(&rock).Error; err != nil {
		t.Fatalf("could not update rock event: %v", err)
	}
	pop := createDAOTestEvent(t, db, "Pop Night", time.Date(2027, 2, 10, 20, 0, 0, 0, time.UTC), 0, false)
	pop.Description = "Pop event"
	pop.Location = "Rosario"
	pop.Price = 1200
	pop.AvailableCapacity = 0
	if err := db.Save(&pop).Error; err != nil {
		t.Fatalf("could not update pop event: %v", err)
	}

	dateFrom := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2027, 1, 31, 23, 59, 59, 0, time.UTC)
	minPrice := 100.0
	maxPrice := 600.0
	isFestival := true
	events, err := eventDAO.GetAll(domain.EventFilters{
		Search:        "rock",
		Location:      "cord",
		DateFrom:      &dateFrom,
		DateTo:        &dateTo,
		MinPrice:      &minPrice,
		MaxPrice:      &maxPrice,
		IsFestival:    &isFestival,
		AvailableOnly: true,
	})
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}

	if len(events) != 1 || events[0].ID != rock.ID {
		t.Fatalf("expected only rock event, got %+v", events)
	}
}

func TestEventDAOGetAllSortsEvents(t *testing.T) {
	db := newDAOTestDB(t)
	eventDAO := NewEventDAO(db)
	first := createDAOTestEvent(t, db, "First", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	first.Price = 300
	second := createDAOTestEvent(t, db, "Second", time.Date(2027, 2, 1, 20, 0, 0, 0, time.UTC), 100, false)
	second.Price = 100
	if err := db.Save(&first).Error; err != nil {
		t.Fatalf("could not save first event: %v", err)
	}
	if err := db.Save(&second).Error; err != nil {
		t.Fatalf("could not save second event: %v", err)
	}

	byDateDesc, err := eventDAO.GetAll(domain.EventFilters{Sort: domain.EventSortDateDesc})
	if err != nil {
		t.Fatalf("GetAll date desc returned error: %v", err)
	}
	if byDateDesc[0].ID != second.ID {
		t.Fatalf("expected second event first by date desc, got %+v", byDateDesc)
	}

	byPriceAsc, err := eventDAO.GetAll(domain.EventFilters{Sort: domain.EventSortPriceAsc})
	if err != nil {
		t.Fatalf("GetAll price asc returned error: %v", err)
	}
	if byPriceAsc[0].ID != second.ID {
		t.Fatalf("expected second event first by price asc, got %+v", byPriceAsc)
	}
}

func TestEventDAOCancelWithCleanupCancelsFestivalAndDeletesAssociatedData(t *testing.T) {
	db := newDAOTestDB(t)
	eventDAO := NewEventDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Festival", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, true)
	schedule := domain.FestivalSchedule{
		EventID:   event.ID,
		Artist:    "Artist",
		Stage:     "Main",
		StartTime: event.Date,
		EndTime:   event.Date.Add(time.Hour),
	}
	if err := db.Create(&schedule).Error; err != nil {
		t.Fatalf("could not create schedule: %v", err)
	}
	itinerary := domain.UserItinerary{UserID: user.ID, EventID: event.ID}
	if err := db.Create(&itinerary).Error; err != nil {
		t.Fatalf("could not create itinerary: %v", err)
	}
	item := domain.ItineraryItem{
		UserItineraryID: itinerary.ID,
		Type:            domain.ItineraryItemTypeShow,
		Title:           "Artist",
		StartTime:       schedule.StartTime,
		EndTime:         schedule.EndTime,
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("could not create itinerary item: %v", err)
	}

	cancelled, err := eventDAO.CancelWithCleanup(event.ID)
	if err != nil {
		t.Fatalf("CancelWithCleanup returned error: %v", err)
	}

	if cancelled.Status != domain.EventStatusCancelled {
		t.Fatalf("expected cancelled status, got %s", cancelled.Status)
	}
	assertCount(t, db, &domain.FestivalSchedule{}, 0)
	assertCount(t, db, &domain.UserItinerary{}, 0)
	assertCount(t, db, &domain.ItineraryItem{}, 0)
}

func TestEventDAOCancelWithCleanupMissingEvent(t *testing.T) {
	db := newDAOTestDB(t)
	eventDAO := NewEventDAO(db)

	_, err := eventDAO.CancelWithCleanup(999)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func createDAOTestEvent(t *testing.T, db *gorm.DB, title string, date time.Time, capacity int, isFestival bool) domain.Event {
	t.Helper()

	event := domain.Event{
		Title:             title,
		Description:       title + " description",
		Date:              date,
		Location:          "Cordoba",
		Capacity:          capacity,
		AvailableCapacity: capacity,
		Price:             100,
		Status:            domain.EventStatusActive,
		IsFestival:        isFestival,
	}
	if err := db.Create(&event).Error; err != nil {
		t.Fatalf("could not create test event: %v", err)
	}

	return event
}

func assertCount(t *testing.T, db *gorm.DB, model interface{}, want int64) {
	t.Helper()

	var got int64
	if err := db.Model(model).Count(&got).Error; err != nil {
		t.Fatalf("could not count model: %v", err)
	}
	if got != want {
		t.Fatalf("expected count %d, got %d", want, got)
	}
}
