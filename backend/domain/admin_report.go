package domain

type AdminSummaryReport struct {
	TotalEvents      int64   `json:"total_events"`
	ActiveEvents     int64   `json:"active_events"`
	CancelledEvents  int64   `json:"cancelled_events"`
	TotalTickets     int64   `json:"total_tickets"`
	ActiveTickets    int64   `json:"active_tickets"`
	CancelledTickets int64   `json:"cancelled_tickets"`
	EstimatedRevenue float64 `json:"estimated_revenue"`
}

type AdminEventReport struct {
	EventID             uint    `json:"event_id"`
	Title               string  `json:"title"`
	Status              string  `json:"status"`
	Capacity            int     `json:"capacity"`
	AvailableCapacity   int     `json:"available_capacity"`
	TicketsIssued       int64   `json:"tickets_issued"`
	ActiveTickets       int64   `json:"active_tickets"`
	CancelledTickets    int64   `json:"cancelled_tickets"`
	OccupancyPercentage float64 `json:"occupancy_percentage"`
	EstimatedRevenue    float64 `json:"estimated_revenue"`
}
