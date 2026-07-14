package dao

import (
	"testing"
	"time"

	"backend/domain"
)

func TestAdminReportDAOGetSummary(t *testing.T) {
	db := newDAOTestDB(t)
	reportDAO := NewAdminReportDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	activeEvent := createDAOTestEvent(t, db, "Active", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	activeEvent.Price = 250
	if err := db.Save(&activeEvent).Error; err != nil {
		t.Fatalf("could not save event: %v", err)
	}
	cancelledEvent := createDAOTestEvent(t, db, "Cancelled", time.Date(2027, 1, 2, 20, 0, 0, 0, time.UTC), 100, false)
	cancelledEvent.Status = domain.EventStatusCancelled
	if err := db.Save(&cancelledEvent).Error; err != nil {
		t.Fatalf("could not save cancelled event: %v", err)
	}
	createDAOTestTicket(t, db, user.ID, activeEvent.ID, "ACTIVE", domain.TicketStatusActive, time.Now())
	createDAOTestTicket(t, db, user.ID, activeEvent.ID, "CANCELLED", domain.TicketStatusCancelled, time.Now())

	report, err := reportDAO.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary returned error: %v", err)
	}

	if report.TotalEvents != 2 || report.ActiveEvents != 1 || report.CancelledEvents != 1 {
		t.Fatalf("unexpected event counts: %+v", report)
	}
	if report.TotalTickets != 2 || report.ActiveTickets != 1 || report.CancelledTickets != 1 {
		t.Fatalf("unexpected ticket counts: %+v", report)
	}
	if report.EstimatedRevenue != 250 {
		t.Fatalf("expected revenue 250, got %f", report.EstimatedRevenue)
	}
}

func TestAdminReportDAOGetEventReports(t *testing.T) {
	db := newDAOTestDB(t)
	reportDAO := NewAdminReportDAO(db)
	user := createDAOTestUser(t, db, "user@mail.com")
	first := createDAOTestEvent(t, db, "First", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, false)
	first.Price = 100
	second := createDAOTestEvent(t, db, "Second", time.Date(2027, 1, 2, 20, 0, 0, 0, time.UTC), 100, false)
	second.Price = 200
	if err := db.Save(&first).Error; err != nil {
		t.Fatalf("could not save first event: %v", err)
	}
	if err := db.Save(&second).Error; err != nil {
		t.Fatalf("could not save second event: %v", err)
	}
	createDAOTestTicket(t, db, user.ID, second.ID, "SECOND-A", domain.TicketStatusActive, time.Now())
	createDAOTestTicket(t, db, user.ID, second.ID, "SECOND-B", domain.TicketStatusActive, time.Now())
	createDAOTestTicket(t, db, user.ID, first.ID, "FIRST-C", domain.TicketStatusCancelled, time.Now())

	reports, err := reportDAO.GetEventReports()
	if err != nil {
		t.Fatalf("GetEventReports returned error: %v", err)
	}

	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
	if reports[0].EventID != second.ID {
		t.Fatalf("expected second event first by active tickets, got %+v", reports)
	}
	if reports[0].ActiveTickets != 2 || reports[0].EstimatedRevenue != 400 {
		t.Fatalf("unexpected first report: %+v", reports[0])
	}
	if reports[1].CancelledTickets != 1 {
		t.Fatalf("expected first event cancelled count 1, got %+v", reports[1])
	}
}
