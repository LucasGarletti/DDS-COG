package services

import "backend/domain"

type AdminReportRepository interface {
	GetSummary() (*domain.AdminSummaryReport, error)
	GetEventReports() ([]domain.AdminEventReport, error)
}

type AdminReportService struct {
	reportDAO AdminReportRepository
}

func NewAdminReportService(reportDAO AdminReportRepository) *AdminReportService {
	return &AdminReportService{reportDAO: reportDAO}
}

func (service *AdminReportService) GetSummary() (*domain.AdminSummaryReport, error) {
	report, err := service.reportDAO.GetSummary()
	if err != nil {
		return nil, err
	}

	if report == nil {
		return &domain.AdminSummaryReport{}, nil
	}

	return report, nil
}

func (service *AdminReportService) GetEventReports() ([]domain.AdminEventReport, error) {
	reports, err := service.reportDAO.GetEventReports()
	if err != nil {
		return nil, err
	}

	if reports == nil {
		return []domain.AdminEventReport{}, nil
	}

	for index := range reports {
		reports[index].OccupancyPercentage = calculateOccupancyPercentage(reports[index].ActiveTickets, reports[index].Capacity)
	}

	return reports, nil
}

func calculateOccupancyPercentage(activeTickets int64, capacity int) float64 {
	if capacity <= 0 {
		return 0
	}

	return float64(activeTickets) / float64(capacity) * 100
}
