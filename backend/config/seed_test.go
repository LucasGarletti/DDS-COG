package config

import (
	"testing"

	"backend/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newSeedTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}, &domain.FestivalSchedule{}, &domain.UserItinerary{}, &domain.ItineraryItem{}); err != nil {
		t.Fatalf("could not migrate test database: %v", err)
	}

	return db
}

func TestSeedDatabaseCreatesUsersWithHashedPasswords(t *testing.T) {
	db := newSeedTestDB(t)

	if err := SeedDatabase(db); err != nil {
		t.Fatalf("SeedDatabase returned error: %v", err)
	}

	var admin domain.User
	if err := db.Where("email = ?", SeedAdminEmail).First(&admin).Error; err != nil {
		t.Fatalf("expected admin user: %v", err)
	}

	if admin.Role != domain.UserRoleAdmin {
		t.Fatalf("expected admin role, got %s", admin.Role)
	}

	if admin.Password == SeedAdminPassword {
		t.Fatal("expected admin password to be hashed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(SeedAdminPassword)); err != nil {
		t.Fatalf("expected admin password hash to match seed password: %v", err)
	}

	var client domain.User
	if err := db.Where("email = ?", SeedClientEmail).First(&client).Error; err != nil {
		t.Fatalf("expected client user: %v", err)
	}

	if client.Role != domain.UserRoleClient {
		t.Fatalf("expected client role, got %s", client.Role)
	}

	if client.Password == SeedClientPassword {
		t.Fatal("expected client password to be hashed")
	}
}

func TestSeedDatabaseCreatesEventsAndSchedules(t *testing.T) {
	db := newSeedTestDB(t)

	if err := SeedDatabase(db); err != nil {
		t.Fatalf("SeedDatabase returned error: %v", err)
	}

	var eventCount int64
	if err := db.Model(&domain.Event{}).Count(&eventCount).Error; err != nil {
		t.Fatalf("could not count events: %v", err)
	}

	if eventCount != 3 {
		t.Fatalf("expected 3 events, got %d", eventCount)
	}

	var cosquinRock domain.Event
	if err := db.Where("title = ?", "Festival Cosquín Rock").First(&cosquinRock).Error; err != nil {
		t.Fatalf("expected Cosquin Rock event: %v", err)
	}

	if !cosquinRock.IsFestival {
		t.Fatal("expected Cosquin Rock to be a festival")
	}

	var scheduleCount int64
	if err := db.Model(&domain.FestivalSchedule{}).Where("event_id = ?", cosquinRock.ID).Count(&scheduleCount).Error; err != nil {
		t.Fatalf("could not count schedules: %v", err)
	}

	if scheduleCount != 3 {
		t.Fatalf("expected 3 schedules, got %d", scheduleCount)
	}
}

func TestSeedDatabaseIsIdempotent(t *testing.T) {
	db := newSeedTestDB(t)

	if err := SeedDatabase(db); err != nil {
		t.Fatalf("first SeedDatabase returned error: %v", err)
	}

	if err := SeedDatabase(db); err != nil {
		t.Fatalf("second SeedDatabase returned error: %v", err)
	}

	counts := []struct {
		name  string
		model interface{}
		want  int64
	}{
		{name: "users", model: &domain.User{}, want: 2},
		{name: "events", model: &domain.Event{}, want: 3},
		{name: "festival schedules", model: &domain.FestivalSchedule{}, want: 3},
	}

	for _, item := range counts {
		var got int64
		if err := db.Model(item.model).Count(&got).Error; err != nil {
			t.Fatalf("could not count %s: %v", item.name, err)
		}

		if got != item.want {
			t.Fatalf("expected %s count %d, got %d", item.name, item.want, got)
		}
	}
}
