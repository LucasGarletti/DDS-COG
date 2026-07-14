package dao

import (
	"testing"

	"backend/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newDAOTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Event{},
		&domain.Ticket{},
		&domain.FestivalSchedule{},
		&domain.UserItinerary{},
		&domain.ItineraryItem{},
	); err != nil {
		t.Fatalf("could not migrate test database: %v", err)
	}

	return db
}

func createDAOTestUser(t *testing.T, db *gorm.DB, email string) domain.User {
	t.Helper()

	user := domain.User{
		Name:     "Test User",
		Email:    email,
		Password: "hashed-password",
		Role:     domain.UserRoleClient,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	return user
}
