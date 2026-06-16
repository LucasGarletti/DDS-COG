package domain

import "time"

const (
	TicketStatusActive    = "active"
	TicketStatusCancelled = "cancelled"
)

type Ticket struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	UserID           uint       `json:"user_id" gorm:"not null;index"`
	EventID          uint       `json:"event_id" gorm:"not null;index"`
	Code             string     `json:"code" gorm:"type:varchar(100);uniqueIndex;not null"`
	Status           string     `json:"status" gorm:"type:varchar(20);not null;default:active"`
	PurchaseDate     time.Time  `json:"purchase_date" gorm:"not null"`
	CancellationDate *time.Time `json:"cancellation_date"`
	User             User       `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Event            Event      `json:"event,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
