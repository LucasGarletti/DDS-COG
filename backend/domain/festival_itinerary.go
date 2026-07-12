package domain

import "time"

const (
	ItineraryItemTypeShow     = "show"
	ItineraryItemTypePersonal = "personal"
)

type FestivalSchedule struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EventID   uint      `json:"event_id" gorm:"not null;index"`
	Artist    string    `json:"artist" gorm:"type:varchar(150);not null"`
	Stage     string    `json:"stage" gorm:"type:varchar(150);not null"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
	ImageURL  string    `json:"image_url" gorm:"type:varchar(255)"`
	Event     Event     `json:"event,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserItinerary struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	UserID    uint            `json:"user_id" gorm:"not null;uniqueIndex:idx_user_event_itinerary"`
	EventID   uint            `json:"event_id" gorm:"not null;uniqueIndex:idx_user_event_itinerary"`
	User      User            `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Event     Event           `json:"event,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Items     []ItineraryItem `json:"items,omitempty" gorm:"foreignKey:UserItineraryID"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ItineraryItem struct {
	ID                 uint              `json:"id" gorm:"primaryKey"`
	UserItineraryID    uint              `json:"user_itinerary_id" gorm:"not null;index"`
	FestivalScheduleID *uint             `json:"festival_schedule_id"`
	Type               string            `json:"type" gorm:"type:varchar(20);not null"`
	Title              string            `json:"title" gorm:"type:varchar(150);not null"`
	Location           string            `json:"location" gorm:"type:varchar(150)"`
	Notes              string            `json:"notes" gorm:"type:text"`
	StartTime          time.Time         `json:"start_time" gorm:"not null"`
	EndTime            time.Time         `json:"end_time" gorm:"not null"`
	UserItinerary      UserItinerary     `json:"user_itinerary,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FestivalSchedule   *FestivalSchedule `json:"festival_schedule,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}
