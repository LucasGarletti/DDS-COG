package dao

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

func TestUserItineraryDAOHasActiveTicket(t *testing.T) {
	db := newDAOTestDB(t)
	itineraryDAO := NewUserItineraryDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)

	hasTicket, err := itineraryDAO.HasActiveTicket(user.ID, event.ID)
	if err != nil {
		t.Fatalf("HasActiveTicket returned error: %v", err)
	}
	if hasTicket {
		t.Fatal("expected no active ticket before creating one")
	}

	createDAOTestTicket(t, db, user.ID, event.ID, "ACTIVE", domain.TicketStatusActive, time.Now())
	hasTicket, err = itineraryDAO.HasActiveTicket(user.ID, event.ID)
	if err != nil {
		t.Fatalf("HasActiveTicket returned error: %v", err)
	}
	if !hasTicket {
		t.Fatal("expected active ticket")
	}
}

func TestUserItineraryDAOCreateAndGetByUserAndEvent(t *testing.T) {
	db := newDAOTestDB(t)
	itineraryDAO := NewUserItineraryDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	itinerary := &domain.UserItinerary{UserID: user.ID, EventID: event.ID}

	if err := itineraryDAO.Create(itinerary); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	found, err := itineraryDAO.GetByUserAndEvent(user.ID, event.ID)
	if err != nil {
		t.Fatalf("GetByUserAndEvent returned error: %v", err)
	}
	if found.ID != itinerary.ID {
		t.Fatalf("expected itinerary %d, got %d", itinerary.ID, found.ID)
	}
}

func TestUserItineraryDAOGetByUserAndEventMissing(t *testing.T) {
	db := newDAOTestDB(t)
	itineraryDAO := NewUserItineraryDAO(db)

	_, err := itineraryDAO.GetByUserAndEvent(1, 1)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func TestUserItineraryDAOCreateGetAndDeleteItems(t *testing.T) {
	db := newDAOTestDB(t)
	itineraryDAO := NewUserItineraryDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	itinerary := &domain.UserItinerary{UserID: user.ID, EventID: event.ID}
	if err := itineraryDAO.Create(itinerary); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	later := &domain.ItineraryItem{
		UserItineraryID: itinerary.ID,
		Type:            domain.ItineraryItemTypePersonal,
		Title:           "Later",
		StartTime:       event.Date.Add(2 * time.Hour),
		EndTime:         event.Date.Add(3 * time.Hour),
	}
	earlier := &domain.ItineraryItem{
		UserItineraryID: itinerary.ID,
		Type:            domain.ItineraryItemTypePersonal,
		Title:           "Earlier",
		StartTime:       event.Date,
		EndTime:         event.Date.Add(time.Hour),
	}

	if err := itineraryDAO.CreateItem(later); err != nil {
		t.Fatalf("CreateItem later returned error: %v", err)
	}
	if err := itineraryDAO.CreateItem(earlier); err != nil {
		t.Fatalf("CreateItem earlier returned error: %v", err)
	}

	items, err := itineraryDAO.GetItemsByItineraryID(itinerary.ID)
	if err != nil {
		t.Fatalf("GetItemsByItineraryID returned error: %v", err)
	}
	if len(items) != 2 || items[0].ID != earlier.ID || items[1].ID != later.ID {
		t.Fatalf("expected items ordered by start time, got %+v", items)
	}

	found, err := itineraryDAO.GetItemByID(earlier.ID)
	if err != nil {
		t.Fatalf("GetItemByID returned error: %v", err)
	}
	if found.Title != "Earlier" {
		t.Fatalf("expected Earlier item, got %+v", found)
	}

	if err := itineraryDAO.DeleteItem(earlier.ID); err != nil {
		t.Fatalf("DeleteItem returned error: %v", err)
	}
	_, err = itineraryDAO.GetItemByID(earlier.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected deleted item not found, got %v", err)
	}
}

func TestUserItineraryDAODeleteByUserAndEvent(t *testing.T) {
	db := newDAOTestDB(t)
	itineraryDAO := NewUserItineraryDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	itinerary := &domain.UserItinerary{UserID: user.ID, EventID: event.ID}
	if err := itineraryDAO.Create(itinerary); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	item := &domain.ItineraryItem{
		UserItineraryID: itinerary.ID,
		Type:            domain.ItineraryItemTypePersonal,
		Title:           "Break",
		StartTime:       event.Date,
		EndTime:         event.Date.Add(time.Hour),
	}
	if err := itineraryDAO.CreateItem(item); err != nil {
		t.Fatalf("CreateItem returned error: %v", err)
	}

	if err := itineraryDAO.DeleteByUserAndEvent(user.ID, event.ID); err != nil {
		t.Fatalf("DeleteByUserAndEvent returned error: %v", err)
	}

	var itineraryCount int64
	if err := db.Model(&domain.UserItinerary{}).Count(&itineraryCount).Error; err != nil {
		t.Fatalf("could not count itineraries: %v", err)
	}
	var itemCount int64
	if err := db.Model(&domain.ItineraryItem{}).Count(&itemCount).Error; err != nil {
		t.Fatalf("could not count items: %v", err)
	}
	if itineraryCount != 0 || itemCount != 0 {
		t.Fatalf("expected itinerary and item deleted, got itineraries=%d items=%d", itineraryCount, itemCount)
	}
}

func TestUserItineraryDAODeleteByUserAndEventNoRows(t *testing.T) {
	db := newDAOTestDB(t)
	itineraryDAO := NewUserItineraryDAO(db)

	if err := itineraryDAO.DeleteByUserAndEvent(1, 1); err != nil {
		t.Fatalf("DeleteByUserAndEvent returned error: %v", err)
	}
}
