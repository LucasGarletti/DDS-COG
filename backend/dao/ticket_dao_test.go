package dao

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

func TestNewTicketDAO(t *testing.T) {
	db := newDAOTestDB(t)

	ticketDAO := NewTicketDAO(db)

	if ticketDAO == nil {
		t.Fatal("expected ticket DAO")
	}
}

func TestTicketDAOGetEventByID(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)

	found, err := ticketDAO.GetEventByID(event.ID)
	if err != nil {
		t.Fatalf("GetEventByID returned error: %v", err)
	}

	if found.ID != event.ID {
		t.Fatalf("expected event %d, got %d", event.ID, found.ID)
	}
}

func TestTicketDAOGetEventByIDMissingEvent(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)

	_, err := ticketDAO.GetEventByID(999)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func TestTicketDAOCreatePurchaseAndCountTickets(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	event.AvailableCapacity = 99
	ticket := &domain.Ticket{
		UserID:       user.ID,
		EventID:      event.ID,
		Code:         "TICKET-2027-000001",
		Status:       domain.TicketStatusActive,
		PurchaseDate: time.Now(),
	}

	if err := ticketDAO.CreatePurchase(ticket, &event); err != nil {
		t.Fatalf("CreatePurchase returned error: %v", err)
	}

	count, err := ticketDAO.CountTickets()
	if err != nil {
		t.Fatalf("CountTickets returned error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 ticket, got %d", count)
	}

	storedEvent, err := ticketDAO.GetEventByID(event.ID)
	if err != nil {
		t.Fatalf("GetEventByID returned error: %v", err)
	}
	if storedEvent.AvailableCapacity != 99 {
		t.Fatalf("expected available capacity 99, got %d", storedEvent.AvailableCapacity)
	}
}

func TestTicketDAOGetByUserIDOrdersAndPreloadsEvent(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	otherUser := createDAOTestUser(t, db, "other@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	oldTicket := createDAOTestTicket(t, db, user.ID, event.ID, "OLD", domain.TicketStatusActive, time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC))
	newTicket := createDAOTestTicket(t, db, user.ID, event.ID, "NEW", domain.TicketStatusActive, time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC))
	createDAOTestTicket(t, db, otherUser.ID, event.ID, "OTHER", domain.TicketStatusActive, time.Date(2027, 1, 3, 12, 0, 0, 0, time.UTC))

	tickets, err := ticketDAO.GetByUserID(user.ID)
	if err != nil {
		t.Fatalf("GetByUserID returned error: %v", err)
	}

	if len(tickets) != 2 {
		t.Fatalf("expected 2 tickets, got %d", len(tickets))
	}
	if tickets[0].ID != newTicket.ID || tickets[1].ID != oldTicket.ID {
		t.Fatalf("expected tickets ordered by purchase date desc, got %+v", tickets)
	}
	if tickets[0].Event.ID != event.ID {
		t.Fatalf("expected event preload, got %+v", tickets[0].Event)
	}
}

func TestTicketDAOGetByIDWithEvent(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	ticket := createDAOTestTicket(t, db, user.ID, event.ID, "CODE", domain.TicketStatusActive, time.Now())

	found, err := ticketDAO.GetByIDWithEvent(ticket.ID)
	if err != nil {
		t.Fatalf("GetByIDWithEvent returned error: %v", err)
	}

	if found.Event.ID != event.ID {
		t.Fatalf("expected event preload, got %+v", found.Event)
	}
}

func TestTicketDAOSaveTicketAndEvent(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	ticket := createDAOTestTicket(t, db, user.ID, event.ID, "CODE", domain.TicketStatusActive, time.Now())
	ticket.Status = domain.TicketStatusCancelled
	event.AvailableCapacity = 100

	if err := ticketDAO.SaveTicketAndEvent(&ticket, &event); err != nil {
		t.Fatalf("SaveTicketAndEvent returned error: %v", err)
	}

	found, err := ticketDAO.GetByIDWithEvent(ticket.ID)
	if err != nil {
		t.Fatalf("GetByIDWithEvent returned error: %v", err)
	}
	if found.Status != domain.TicketStatusCancelled || found.Event.AvailableCapacity != 100 {
		t.Fatalf("expected saved ticket and event, got %+v", found)
	}
}

func TestTicketDAOSaveTicket(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	ticket := createDAOTestTicket(t, db, user.ID, event.ID, "CODE", domain.TicketStatusActive, time.Now())
	ticket.UserID = user.ID
	ticket.Status = domain.TicketStatusCancelled

	if err := ticketDAO.SaveTicket(&ticket); err != nil {
		t.Fatalf("SaveTicket returned error: %v", err)
	}

	found, err := ticketDAO.GetByIDWithEvent(ticket.ID)
	if err != nil {
		t.Fatalf("GetByIDWithEvent returned error: %v", err)
	}
	if found.Status != domain.TicketStatusCancelled {
		t.Fatalf("expected cancelled ticket, got %s", found.Status)
	}
}

func TestTicketDAOCountsByEventAndStatus(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	otherEvent := createDAOTestEvent(t, db, "Other", time.Date(2027, 1, 2, 20, 0, 0, 0, time.UTC), 100, false)
	createDAOTestTicket(t, db, user.ID, event.ID, "ACTIVE", domain.TicketStatusActive, time.Now())
	createDAOTestTicket(t, db, user.ID, event.ID, "CANCELLED", domain.TicketStatusCancelled, time.Now())
	createDAOTestTicket(t, db, user.ID, otherEvent.ID, "OTHER", domain.TicketStatusActive, time.Now())

	total, err := ticketDAO.CountTicketsByEvent(event.ID)
	if err != nil {
		t.Fatalf("CountTicketsByEvent returned error: %v", err)
	}
	active, err := ticketDAO.CountTicketsByEventAndStatus(event.ID, domain.TicketStatusActive)
	if err != nil {
		t.Fatalf("CountTicketsByEventAndStatus returned error: %v", err)
	}
	cancelled, err := ticketDAO.CountTicketsByEventAndStatus(event.ID, domain.TicketStatusCancelled)
	if err != nil {
		t.Fatalf("CountTicketsByEventAndStatus returned error: %v", err)
	}

	if total != 2 || active != 1 || cancelled != 1 {
		t.Fatalf("unexpected counts total=%d active=%d cancelled=%d", total, active, cancelled)
	}
}

func TestTicketDAODeleteUserItineraryByUserAndEvent(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	otherUser := createDAOTestUser(t, db, "other@mail.com")
	event := createDAOTestEvent(t, db, "Event", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, true)
	itinerary := domain.UserItinerary{UserID: user.ID, EventID: event.ID}
	if err := db.Create(&itinerary).Error; err != nil {
		t.Fatalf("could not create itinerary: %v", err)
	}
	item := domain.ItineraryItem{
		UserItineraryID: itinerary.ID,
		Type:            domain.ItineraryItemTypePersonal,
		Title:           "Break",
		StartTime:       event.Date,
		EndTime:         event.Date.Add(time.Hour),
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("could not create itinerary item: %v", err)
	}
	otherItinerary := domain.UserItinerary{UserID: otherUser.ID, EventID: event.ID}
	if err := db.Create(&otherItinerary).Error; err != nil {
		t.Fatalf("could not create other itinerary: %v", err)
	}

	if err := ticketDAO.DeleteUserItineraryByUserAndEvent(user.ID, event.ID); err != nil {
		t.Fatalf("DeleteUserItineraryByUserAndEvent returned error: %v", err)
	}

	var userItineraries int64
	if err := db.Model(&domain.UserItinerary{}).Where("user_id = ?", user.ID).Count(&userItineraries).Error; err != nil {
		t.Fatalf("could not count user itineraries: %v", err)
	}
	if userItineraries != 0 {
		t.Fatalf("expected user itinerary deleted, got %d", userItineraries)
	}
	var items int64
	if err := db.Model(&domain.ItineraryItem{}).Where("user_itinerary_id = ?", itinerary.ID).Count(&items).Error; err != nil {
		t.Fatalf("could not count itinerary items: %v", err)
	}
	if items != 0 {
		t.Fatalf("expected itinerary items deleted, got %d", items)
	}
	var otherItineraries int64
	if err := db.Model(&domain.UserItinerary{}).Where("user_id = ?", otherUser.ID).Count(&otherItineraries).Error; err != nil {
		t.Fatalf("could not count other itineraries: %v", err)
	}
	if otherItineraries != 1 {
		t.Fatalf("expected other user itinerary to remain, got %d", otherItineraries)
	}
}

func TestTicketDAODeleteUserItineraryByUserAndEventNoRows(t *testing.T) {
	db := newDAOTestDB(t)
	ticketDAO := NewTicketDAO(db)

	if err := ticketDAO.DeleteUserItineraryByUserAndEvent(1, 1); err != nil {
		t.Fatalf("DeleteUserItineraryByUserAndEvent returned error: %v", err)
	}
}

func createDAOTestTicket(t *testing.T, db *gorm.DB, userID uint, eventID uint, code string, status string, purchaseDate time.Time) domain.Ticket {
	t.Helper()

	ticket := domain.Ticket{
		UserID:       userID,
		EventID:      eventID,
		Code:         code,
		Status:       status,
		PurchaseDate: purchaseDate,
	}
	if err := db.Create(&ticket).Error; err != nil {
		t.Fatalf("could not create test ticket: %v", err)
	}

	return ticket
}
