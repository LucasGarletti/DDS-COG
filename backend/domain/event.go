package domain

import "time"

type Event struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Title             string    `json:"title" gorm:"type:varchar(150);not null"`
	Description       string    `json:"description" gorm:"type:text"`
	Date              time.Time `json:"date" gorm:"not null"`
	Location          string    `json:"location" gorm:"type:varchar(150);not null"`
	Capacity          int       `json:"capacity" gorm:"not null"`
	AvailableCapacity int       `json:"available_capacity" gorm:"not null"`
	Price             float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	ImageURL          string    `json:"image_url" gorm:"type:varchar(255)"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
