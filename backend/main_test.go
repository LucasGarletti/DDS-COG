package main

import (
	"testing"
	"time"

	"backend/config"
	"backend/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newMainTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Event{}); err != nil {
		t.Fatalf("could not migrate test database: %v", err)
	}

	return db
}

func TestNormalizeUserRoles(t *testing.T) {
	originalDB := config.DB
	t.Cleanup(func() { config.DB = originalDB })
	config.DB = newMainTestDB(t)

	users := []domain.User{
		{Name: "Valid Client", Email: "client@mail.com", Password: "password", Role: domain.UserRoleClient},
		{Name: "Valid Admin", Email: "admin@mail.com", Password: "password", Role: domain.UserRoleAdmin},
		{Name: "Empty Role", Email: "empty@mail.com", Password: "password", Role: ""},
		{Name: "Invalid Role", Email: "invalid@mail.com", Password: "password", Role: "superadmin"},
	}
	for _, user := range users {
		if err := config.DB.Create(&user).Error; err != nil {
			t.Fatalf("could not create user: %v", err)
		}
	}

	if err := normalizeUserRoles(); err != nil {
		t.Fatalf("normalizeUserRoles returned error: %v", err)
	}

	var normalized domain.User
	if err := config.DB.Where("email = ?", "invalid@mail.com").First(&normalized).Error; err != nil {
		t.Fatalf("could not find invalid user: %v", err)
	}
	if normalized.Role != domain.UserRoleClient {
		t.Fatalf("expected invalid role normalized to client, got %s", normalized.Role)
	}

	var admin domain.User
	if err := config.DB.Where("email = ?", "admin@mail.com").First(&admin).Error; err != nil {
		t.Fatalf("could not find admin: %v", err)
	}
	if admin.Role != domain.UserRoleAdmin {
		t.Fatalf("expected admin role to remain admin, got %s", admin.Role)
	}
}

func TestNormalizeEventStatuses(t *testing.T) {
	originalDB := config.DB
	t.Cleanup(func() { config.DB = originalDB })
	config.DB = newMainTestDB(t)

	events := []domain.Event{
		{
			Title:             "Active",
			Date:              time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC),
			Location:          "Cordoba",
			Capacity:          100,
			AvailableCapacity: 100,
			Price:             100,
			Status:            domain.EventStatusActive,
		},
		{
			Title:             "Cancelled",
			Date:              time.Date(2027, 1, 2, 20, 0, 0, 0, time.UTC),
			Location:          "Cordoba",
			Capacity:          100,
			AvailableCapacity: 100,
			Price:             100,
			Status:            domain.EventStatusCancelled,
		},
		{
			Title:             "Invalid",
			Date:              time.Date(2027, 1, 3, 20, 0, 0, 0, time.UTC),
			Location:          "Cordoba",
			Capacity:          100,
			AvailableCapacity: 100,
			Price:             100,
			Status:            "archived",
		},
	}
	for _, event := range events {
		if err := config.DB.Create(&event).Error; err != nil {
			t.Fatalf("could not create event: %v", err)
		}
	}

	if err := normalizeEventStatuses(); err != nil {
		t.Fatalf("normalizeEventStatuses returned error: %v", err)
	}

	var normalized domain.Event
	if err := config.DB.Where("title = ?", "Invalid").First(&normalized).Error; err != nil {
		t.Fatalf("could not find invalid event: %v", err)
	}
	if normalized.Status != domain.EventStatusActive {
		t.Fatalf("expected invalid status normalized to active, got %s", normalized.Status)
	}

	var cancelled domain.Event
	if err := config.DB.Where("title = ?", "Cancelled").First(&cancelled).Error; err != nil {
		t.Fatalf("could not find cancelled event: %v", err)
	}
	if cancelled.Status != domain.EventStatusCancelled {
		t.Fatalf("expected cancelled status to remain cancelled, got %s", cancelled.Status)
	}
}
