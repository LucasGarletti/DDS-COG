package domain

import "time"

const (
	EventSortDateAsc   = "date_asc"
	EventSortDateDesc  = "date_desc"
	EventSortPriceAsc  = "price_asc"
	EventSortPriceDesc = "price_desc"
)

type EventFilters struct {
	Search        string
	Location      string
	DateFrom      *time.Time
	DateTo        *time.Time
	MinPrice      *float64
	MaxPrice      *float64
	IsFestival    *bool
	AvailableOnly bool
	Sort          string
}
