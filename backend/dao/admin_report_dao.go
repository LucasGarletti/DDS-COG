package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type AdminReportDAO struct {
	db *gorm.DB
}

func NewAdminReportDAO(db *gorm.DB) *AdminReportDAO {
	return &AdminReportDAO{db: db}
}

func (dao *AdminReportDAO) GetSummary() (*domain.AdminSummaryReport, error) {
	var report domain.AdminSummaryReport

	err := dao.db.Raw(`
		SELECT
			(SELECT COUNT(*) FROM events) AS total_events,
			(SELECT COUNT(*) FROM events WHERE status = ?) AS active_events,
			(SELECT COUNT(*) FROM events WHERE status = ?) AS cancelled_events,
			(SELECT COUNT(*) FROM tickets) AS total_tickets,
			(SELECT COUNT(*) FROM tickets WHERE status = ?) AS active_tickets,
			(SELECT COUNT(*) FROM tickets WHERE status = ?) AS cancelled_tickets,
			COALESCE((
				SELECT SUM(events.price)
				FROM tickets
				JOIN events ON events.id = tickets.event_id
				WHERE tickets.status = ?
			), 0) AS estimated_revenue
	`,
		domain.EventStatusActive,
		domain.EventStatusCancelled,
		domain.TicketStatusActive,
		domain.TicketStatusCancelled,
		domain.TicketStatusActive,
	).Scan(&report).Error
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (dao *AdminReportDAO) GetEventReports() ([]domain.AdminEventReport, error) {
	var reports []domain.AdminEventReport

	err := dao.db.Raw(`
		SELECT
			events.id AS event_id,
			events.title AS title,
			events.status AS status,
			events.capacity AS capacity,
			events.available_capacity AS available_capacity,
			COUNT(tickets.id) AS tickets_issued,
			COALESCE(SUM(CASE WHEN tickets.status = ? THEN 1 ELSE 0 END), 0) AS active_tickets,
			COALESCE(SUM(CASE WHEN tickets.status = ? THEN 1 ELSE 0 END), 0) AS cancelled_tickets,
			COALESCE(SUM(CASE WHEN tickets.status = ? THEN events.price ELSE 0 END), 0) AS estimated_revenue
		FROM events
		LEFT JOIN tickets ON tickets.event_id = events.id
		GROUP BY events.id, events.title, events.status, events.capacity, events.available_capacity
		ORDER BY active_tickets DESC, events.id ASC
	`,
		domain.TicketStatusActive,
		domain.TicketStatusCancelled,
		domain.TicketStatusActive,
	).Scan(&reports).Error
	if err != nil {
		return nil, err
	}

	return reports, nil
}
