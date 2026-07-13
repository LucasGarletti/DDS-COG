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
	saved     *domain.FestivalSchedule
	deletedID uint
}

func (repo *fakeFestivalScheduleRepository) Create(schedule *domain.FestivalSchedule) error {
	repo.created = schedule
	return repo.err
}

func (repo *fakeFestivalScheduleRepository) Save(schedule *domain.FestivalSchedule) error {
	repo.saved = schedule
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

func validUpdateScheduleInput() UpdateFestivalScheduleInput {
	return UpdateFestivalScheduleInput{
		ID:        3,
		Artist:    "Artista editado",
		Stage:     "Escenario",
		StartTime: time.Date(2026, 12, 10, 20, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2026, 12, 10, 21, 0, 0, 0, time.UTC),
	}
}

func festivalEvent() *domain.Event {
	return &domain.Event{
		ID:         1,
		Date:       time.Date(2026, 12, 10, 16, 0, 0, 0, time.UTC),
		IsFestival: true,
	}
}

func TestFestivalScheduleCreateValid(t *testing.T) {
	repo := &fakeFestivalScheduleRepository{}
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, repo)

	schedule, err := service.Create(validScheduleInput())
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if schedule.Artist != "Artista" || repo.created == nil {
		t.Fatalf("expected schedule to be created, got %+v", schedule)
	}
}

func TestFestivalScheduleCreateInvalidArtist(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, &fakeFestivalScheduleRepository{})
	input := validScheduleInput()
	input.Artist = " "

	_, err := service.Create(input)
	if !errors.Is(err, ErrInvalidFestivalSchedule) {
		t.Fatalf("expected ErrInvalidFestivalSchedule, got %v", err)
	}
}

func TestFestivalScheduleCreateEndBeforeStart(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, &fakeFestivalScheduleRepository{})
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
	event := festivalEvent()
	event.IsFestival = false
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: event}, &fakeFestivalScheduleRepository{})

	_, err := service.Create(validScheduleInput())
	if !errors.Is(err, ErrEventNotFestival) {
		t.Fatalf("expected ErrEventNotFestival, got %v", err)
	}
}

func TestFestivalScheduleListReturnsEmptySlice(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, &fakeFestivalScheduleRepository{})

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

func TestFestivalScheduleCreatePreviousDayReturnsDateMismatch(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, &fakeFestivalScheduleRepository{})
	input := validScheduleInput()
	input.StartTime = time.Date(2026, 12, 9, 20, 0, 0, 0, time.UTC)
	input.EndTime = time.Date(2026, 12, 9, 21, 0, 0, 0, time.UTC)

	_, err := service.Create(input)
	if !errors.Is(err, ErrFestivalScheduleDateMismatch) {
		t.Fatalf("expected ErrFestivalScheduleDateMismatch, got %v", err)
	}
}

func TestFestivalScheduleCreateNextDayReturnsDateMismatch(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, &fakeFestivalScheduleRepository{})
	input := validScheduleInput()
	input.StartTime = time.Date(2026, 12, 11, 20, 0, 0, 0, time.UTC)
	input.EndTime = time.Date(2026, 12, 11, 21, 0, 0, 0, time.UTC)

	_, err := service.Create(input)
	if !errors.Is(err, ErrFestivalScheduleDateMismatch) {
		t.Fatalf("expected ErrFestivalScheduleDateMismatch, got %v", err)
	}
}

func TestFestivalScheduleCreateEndOnNextDayReturnsDateMismatch(t *testing.T) {
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, &fakeFestivalScheduleRepository{})
	input := validScheduleInput()
	input.StartTime = time.Date(2026, 12, 10, 23, 30, 0, 0, time.UTC)
	input.EndTime = time.Date(2026, 12, 11, 0, 30, 0, 0, time.UTC)

	_, err := service.Create(input)
	if !errors.Is(err, ErrFestivalScheduleDateMismatch) {
		t.Fatalf("expected ErrFestivalScheduleDateMismatch, got %v", err)
	}
}

func TestFestivalScheduleUpdateWrongDayReturnsDateMismatch(t *testing.T) {
	repo := &fakeFestivalScheduleRepository{schedule: &domain.FestivalSchedule{ID: 3, EventID: 1}}
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, repo)
	input := validUpdateScheduleInput()
	input.StartTime = time.Date(2026, 12, 11, 20, 0, 0, 0, time.UTC)
	input.EndTime = time.Date(2026, 12, 11, 21, 0, 0, 0, time.UTC)

	_, err := service.Update(input)
	if !errors.Is(err, ErrFestivalScheduleDateMismatch) {
		t.Fatalf("expected ErrFestivalScheduleDateMismatch, got %v", err)
	}
}

func TestFestivalScheduleUpdateValid(t *testing.T) {
	repo := &fakeFestivalScheduleRepository{schedule: &domain.FestivalSchedule{ID: 3, EventID: 1}}
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: festivalEvent()}, repo)

	schedule, err := service.Update(validUpdateScheduleInput())
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if schedule.Artist != "Artista editado" || repo.saved == nil {
		t.Fatalf("expected schedule to be saved, got %+v", schedule)
	}
}

func TestFestivalScheduleDateValidationUsesEventLocation(t *testing.T) {
	location := time.FixedZone("ART", -3*60*60)
	event := &domain.Event{
		ID:         1,
		Date:       time.Date(2027, 2, 15, 16, 0, 0, 0, location),
		IsFestival: true,
	}
	service := NewFestivalScheduleService(fakeFestivalEventRepository{event: event}, &fakeFestivalScheduleRepository{})
	input := CreateFestivalScheduleInput{
		EventID:   1,
		Artist:    "Artista",
		Stage:     "Escenario",
		StartTime: time.Date(2027, 2, 16, 2, 30, 0, 0, time.UTC),
		EndTime:   time.Date(2027, 2, 16, 2, 45, 0, 0, time.UTC),
	}

	schedule, err := service.Create(input)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if schedule == nil {
		t.Fatal("expected schedule")
	}
}
