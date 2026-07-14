package dao

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

func TestFestivalScheduleDAOCreateSaveGetAndDelete(t *testing.T) {
	db := newDAOTestDB(t)
	scheduleDAO := NewFestivalScheduleDAO(db)
	event := createDAOTestEvent(t, db, "Festival", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, true)
	schedule := &domain.FestivalSchedule{
		EventID:   event.ID,
		Artist:    "Artist",
		Stage:     "Main",
		StartTime: event.Date,
		EndTime:   event.Date.Add(time.Hour),
	}

	if err := scheduleDAO.Create(schedule); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	schedule.Stage = "Second"
	if err := scheduleDAO.Save(schedule); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	found, err := scheduleDAO.GetByID(schedule.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if found.Stage != "Second" {
		t.Fatalf("expected updated stage, got %s", found.Stage)
	}

	if err := scheduleDAO.Delete(schedule.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	_, err = scheduleDAO.GetByID(schedule.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func TestFestivalScheduleDAOGetByEventIDOrdersByStartTime(t *testing.T) {
	db := newDAOTestDB(t)
	scheduleDAO := NewFestivalScheduleDAO(db)
	event := createDAOTestEvent(t, db, "Festival", time.Date(2027, 1, 1, 20, 0, 0, 0, time.UTC), 100, true)
	later := domain.FestivalSchedule{EventID: event.ID, Artist: "Later", Stage: "Main", StartTime: event.Date.Add(2 * time.Hour), EndTime: event.Date.Add(3 * time.Hour)}
	earlier := domain.FestivalSchedule{EventID: event.ID, Artist: "Earlier", Stage: "Main", StartTime: event.Date, EndTime: event.Date.Add(time.Hour)}
	if err := db.Create(&later).Error; err != nil {
		t.Fatalf("could not create later schedule: %v", err)
	}
	if err := db.Create(&earlier).Error; err != nil {
		t.Fatalf("could not create earlier schedule: %v", err)
	}

	schedules, err := scheduleDAO.GetByEventID(event.ID)
	if err != nil {
		t.Fatalf("GetByEventID returned error: %v", err)
	}

	if len(schedules) != 2 || schedules[0].ID != earlier.ID || schedules[1].ID != later.ID {
		t.Fatalf("expected schedules ordered by start time, got %+v", schedules)
	}
}
