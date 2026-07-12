package services

import (
	"errors"
	"testing"

	"backend/domain"
)

type fakeAdminReportRepository struct {
	summary      *domain.AdminSummaryReport
	eventReports []domain.AdminEventReport
	err          error
}

func (repo fakeAdminReportRepository) GetSummary() (*domain.AdminSummaryReport, error) {
	if repo.err != nil {
		return nil, repo.err
	}

	return repo.summary, nil
}

func (repo fakeAdminReportRepository) GetEventReports() ([]domain.AdminEventReport, error) {
	if repo.err != nil {
		return nil, repo.err
	}

	return repo.eventReports, nil
}

func TestAdminReportSummaryWithData(t *testing.T) {
	expected := &domain.AdminSummaryReport{
		TotalEvents:      3,
		ActiveEvents:     2,
		CancelledEvents:  1,
		TotalTickets:     10,
		ActiveTickets:    8,
		CancelledTickets: 2,
		EstimatedRevenue: 120000,
	}
	service := NewAdminReportService(fakeAdminReportRepository{summary: expected})

	report, err := service.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary returned error: %v", err)
	}

	if report.TotalEvents != expected.TotalEvents || report.EstimatedRevenue != expected.EstimatedRevenue {
		t.Fatalf("unexpected summary report: %+v", report)
	}
}

func TestAdminReportSummaryWithoutDataReturnsZeroValues(t *testing.T) {
	service := NewAdminReportService(fakeAdminReportRepository{})

	report, err := service.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary returned error: %v", err)
	}

	if report.TotalEvents != 0 || report.TotalTickets != 0 || report.EstimatedRevenue != 0 {
		t.Fatalf("expected zero-value summary, got %+v", report)
	}
}

func TestAdminReportSummaryRepositoryError(t *testing.T) {
	repoErr := errors.New("repo error")
	service := NewAdminReportService(fakeAdminReportRepository{err: repoErr})

	_, err := service.GetSummary()
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestAdminReportEventRankingWithSeveralEvents(t *testing.T) {
	service := NewAdminReportService(fakeAdminReportRepository{
		eventReports: []domain.AdminEventReport{
			{EventID: 2, Title: "B", Capacity: 100, ActiveTickets: 80, TicketsIssued: 90, CancelledTickets: 10, EstimatedRevenue: 8000},
			{EventID: 1, Title: "A", Capacity: 100, ActiveTickets: 40, TicketsIssued: 50, CancelledTickets: 10, EstimatedRevenue: 4000},
		},
	})

	reports, err := service.GetEventReports()
	if err != nil {
		t.Fatalf("GetEventReports returned error: %v", err)
	}

	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}

	if reports[0].EventID != 2 || reports[1].EventID != 1 {
		t.Fatalf("expected repository ranking order to be preserved, got %+v", reports)
	}
}

func TestAdminReportEventOccupancyCalculation(t *testing.T) {
	service := NewAdminReportService(fakeAdminReportRepository{
		eventReports: []domain.AdminEventReport{{EventID: 1, Capacity: 200, ActiveTickets: 50}},
	})

	reports, err := service.GetEventReports()
	if err != nil {
		t.Fatalf("GetEventReports returned error: %v", err)
	}

	if reports[0].OccupancyPercentage != 25 {
		t.Fatalf("expected occupancy 25, got %f", reports[0].OccupancyPercentage)
	}
}

func TestAdminReportEventZeroCapacityOccupancyIsZero(t *testing.T) {
	service := NewAdminReportService(fakeAdminReportRepository{
		eventReports: []domain.AdminEventReport{{EventID: 1, Capacity: 0, ActiveTickets: 50}},
	})

	reports, err := service.GetEventReports()
	if err != nil {
		t.Fatalf("GetEventReports returned error: %v", err)
	}

	if reports[0].OccupancyPercentage != 0 {
		t.Fatalf("expected occupancy 0, got %f", reports[0].OccupancyPercentage)
	}
}

func TestAdminReportEventRankingTieOrderByEventID(t *testing.T) {
	service := NewAdminReportService(fakeAdminReportRepository{
		eventReports: []domain.AdminEventReport{
			{EventID: 1, ActiveTickets: 10},
			{EventID: 2, ActiveTickets: 10},
		},
	})

	reports, err := service.GetEventReports()
	if err != nil {
		t.Fatalf("GetEventReports returned error: %v", err)
	}

	if reports[0].EventID != 1 || reports[1].EventID != 2 {
		t.Fatalf("expected repository tie order to be preserved, got %+v", reports)
	}
}

func TestAdminReportEventWithoutTickets(t *testing.T) {
	service := NewAdminReportService(fakeAdminReportRepository{
		eventReports: []domain.AdminEventReport{{EventID: 1, Title: "Sin tickets", Capacity: 100}},
	})

	reports, err := service.GetEventReports()
	if err != nil {
		t.Fatalf("GetEventReports returned error: %v", err)
	}

	if reports[0].TicketsIssued != 0 || reports[0].ActiveTickets != 0 || reports[0].EstimatedRevenue != 0 {
		t.Fatalf("expected zero ticket fields, got %+v", reports[0])
	}
}

func TestAdminReportEventReportsRepositoryError(t *testing.T) {
	repoErr := errors.New("repo error")
	service := NewAdminReportService(fakeAdminReportRepository{err: repoErr})

	_, err := service.GetEventReports()
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
