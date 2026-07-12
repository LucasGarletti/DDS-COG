package services

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

type fakeFestivalEventRepository struct {
	event *domain.Event
	err   error
}

func (repo fakeFestivalEventRepository) GetByIDForAdmin(id uint) (*domain.Event, error) {
	if repo.err != nil {
		return nil, repo.err
	}
	return repo.event, nil
}

type fakeFestivalScheduleRepository struct {
	schedule  *domain.FestivalSchedule
	schedules []domain.FestivalSchedule
	err       error
	created   *domain.FestivalSchedule
	deletedID uint
}

func (repo *fakeFestivalScheduleRepository) Create(schedule *domain.FestivalSchedule) error {
	repo.created = schedule
	return repo.err
}

func (repo *fakeFestivalScheduleRepository) GetByID(id uint) (*domain.FestivalSchedule, error) {
	if repo.err != nil {
		return nil, repo.err
	}
	return repo.schedule, nil
}

func (repo *fakeFestivalScheduleRepository) GetByEventID(eventID uint) ([]domain.FestivalSchedule, error) {
	return repo.schedules, repo.err
}

func (repo *fakeFestivalScheduleRepository) Delete(id uint) error {
	repo.deletedID = id
	return repo.err
}

func validScheduleInput() CreateFestivalScheduleInput {
	return CreateFestivalScheduleInput{
		EventID:   1,
		Artist:    "Artista",
		Stage:     "Escenario",
		StartTime: time.Date(2026, 12, 10, 20, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 12, 10, 21, 0, 0, 0, time.UTC),
	}
}

func TestFestivalScheduleCreateValid(t *testing.T) {
	repo := &fakeFestivalScheduleRepository{}
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: &domain.Event{ID: 1, IsFestival: true}}, repo)

	schedule, err := service.Create(validScheduleInput())
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if schedule.Artist != "Artista" || repo.created == nil {
		t.Fatalf("expected schedule to be created, got %+v", schedule)
	}
}

func TestFestivalScheduleCreateInvalidArtist(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: &domain.Event{ID: 1, IsFestival: true}}, &fakeFestivalScheduleRepository{})
	input := validScheduleInput()
	input.Artist = " "

	_, err := service.Create(input)
	if !errors.Is(err, ErrInvalidFestivalSchedule) {
		t.Fatalf("expected ErrInvalidFestivalSchedule, got %v", err)
	}
}

func TestFestivalScheduleCreateEndBeforeStart(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: &domain.Event{ID: 1, IsFestival: true}}, &fakeFestivalScheduleRepository{})
	input := validScheduleInput()
	input.EndTime = input.StartTime

	_, err := service.Create(input)
	if !errors.Is(err, ErrInvalidFestivalSchedule) {
		t.Fatalf("expected ErrInvalidFestivalSchedule, got %v", err)
	}
}

func TestFestivalScheduleCreateMissingEventReturnsError(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{err: gorm.ErrRecordNotFound}, &fakeFestivalScheduleRepository{})

	_, err := service.Create(validScheduleInput())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestFestivalScheduleCreateNonFestivalReturnsError(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: &domain.Event{ID: 1, IsFestival: false}}, &fakeFestivalScheduleRepository{})

	_, err := service.Create(validScheduleInput())
	if !errors.Is(err, ErrEventNotFestival) {
		t.Fatalf("expected ErrEventNotFestival, got %v", err)
	}
}

func TestFestivalScheduleListReturnsEmptySlice(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: &domain.Event{ID: 1, IsFestival: true}}, &fakeFestivalScheduleRepository{})

	schedules, err := service.List(1)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if schedules == nil || len(schedules) != 0 {
		t.Fatalf("expected empty schedule slice, got %+v", schedules)
	}
}

func TestFestivalScheduleDeleteSuccess(t *testing.T) {
	repo := &fakeFestivalScheduleRepository{schedule: &domain.FestivalSchedule{ID: 1}}
	service := NewFestivalScheduleService(fakeFestivalEventRepository{}, repo)

	if err := service.Delete(1); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if repo.deletedID != 1 {
		t.Fatalf("expected deleted id 1, got %d", repo.deletedID)
	}
}
