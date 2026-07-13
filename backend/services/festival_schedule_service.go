package services

import (
	"errors"
	"strings"
	"time"

	"backend/domain"
)

var (
	ErrInvalidFestivalSchedule      = errors.New("invalid festival schedule")
	ErrEventNotFestival             = errors.New("event is not festival")
	ErrFestivalScheduleDateMismatch = errors.New("festival schedule date mismatch")
)

type FestivalEventRepository interface {
	GetByIDForAdmin(id uint) (*domain.Event, error)
}

type FestivalScheduleRepository interface {
	Create(schedule *domain.FestivalSchedule) error
	Save(schedule *domain.FestivalSchedule) error
	GetByID(id uint) (*domain.FestivalSchedule, error)
	GetByEventID(eventID uint) ([]domain.FestivalSchedule, error)
	Delete(id uint) error
}

type FestivalScheduleService struct {
	eventDAO    FestivalEventRepository
	scheduleDAO FestivalScheduleRepository
}

type CreateFestivalScheduleInput struct {
	EventID   uint
	Artist    string
	Stage     string
	StartTime time.Time
	EndTime   time.Time
	ImageURL  string
}

type UpdateFestivalScheduleInput struct {
	ID        uint
	Artist    string
	Stage     string
	StartTime time.Time
	EndTime   time.Time
	ImageURL  string
}

type FestivalScheduleDateMismatchError struct {
	FestivalDate time.Time
}

func (err FestivalScheduleDateMismatchError) Error() string {
	return ErrFestivalScheduleDateMismatch.Error()
}

func (err FestivalScheduleDateMismatchError) Is(target error) bool {
	return target == ErrFestivalScheduleDateMismatch
}

func NewFestivalScheduleService(eventDAO FestivalEventRepository, scheduleDAO FestivalScheduleRepository) *FestivalScheduleService {
	return &FestivalScheduleService{
		eventDAO:    eventDAO,
		scheduleDAO: scheduleDAO,
	}
}

func (service *FestivalScheduleService) Create(input CreateFestivalScheduleInput) (*domain.FestivalSchedule, error) {
	input.Artist = strings.TrimSpace(input.Artist)
	input.Stage = strings.TrimSpace(input.Stage)

	if input.EventID == 0 || input.Artist == "" || input.Stage == "" || input.StartTime.IsZero() || input.EndTime.IsZero() || !input.EndTime.After(input.StartTime) {
		return nil, ErrInvalidFestivalSchedule
	}

	event, err := service.eventDAO.GetByIDForAdmin(input.EventID)
	if err != nil {
		return nil, err
	}

	if !event.IsFestival {
		return nil, ErrEventNotFestival
	}

	if err := validateScheduleDateMatchesEvent(event, input.StartTime, input.EndTime); err != nil {
		return nil, err
	}

	schedule := &domain.FestivalSchedule{
		EventID:   input.EventID,
		Artist:    input.Artist,
		Stage:     input.Stage,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		ImageURL:  strings.TrimSpace(input.ImageURL),
	}

	if err := service.scheduleDAO.Create(schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (service *FestivalScheduleService) Update(input UpdateFestivalScheduleInput) (*domain.FestivalSchedule, error) {
	input.Artist = strings.TrimSpace(input.Artist)
	input.Stage = strings.TrimSpace(input.Stage)

	if input.ID == 0 || input.Artist == "" || input.Stage == "" || input.StartTime.IsZero() || input.EndTime.IsZero() || !input.EndTime.After(input.StartTime) {
		return nil, ErrInvalidFestivalSchedule
	}

	schedule, err := service.scheduleDAO.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	event, err := service.eventDAO.GetByIDForAdmin(schedule.EventID)
	if err != nil {
		return nil, err
	}

	if !event.IsFestival {
		return nil, ErrEventNotFestival
	}

	if err := validateScheduleDateMatchesEvent(event, input.StartTime, input.EndTime); err != nil {
		return nil, err
	}

	schedule.Artist = input.Artist
	schedule.Stage = input.Stage
	schedule.StartTime = input.StartTime
	schedule.EndTime = input.EndTime
	schedule.ImageURL = strings.TrimSpace(input.ImageURL)

	if err := service.scheduleDAO.Save(schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (service *FestivalScheduleService) List(eventID uint) ([]domain.FestivalSchedule, error) {
	if eventID == 0 {
		return nil, ErrInvalidFestivalSchedule
	}

	event, err := service.eventDAO.GetByIDForAdmin(eventID)
	if err != nil {
		return nil, err
	}

	if !event.IsFestival {
		return nil, ErrEventNotFestival
	}

	schedules, err := service.scheduleDAO.GetByEventID(eventID)
	if err != nil {
		return nil, err
	}

	if schedules == nil {
		return []domain.FestivalSchedule{}, nil
	}

	return schedules, nil
}

func (service *FestivalScheduleService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidFestivalSchedule
	}

	if _, err := service.scheduleDAO.GetByID(id); err != nil {
		return err
	}

	return service.scheduleDAO.Delete(id)
}

func validateScheduleDateMatchesEvent(event *domain.Event, startTime time.Time, endTime time.Time) error {
	location := event.Date.Location()
	if location == nil {
		location = time.Local
	}

	if !sameCalendarDate(startTime.In(location), event.Date.In(location)) || !sameCalendarDate(endTime.In(location), event.Date.In(location)) {
		return FestivalScheduleDateMismatchError{FestivalDate: event.Date}
	}

	return nil
}

func sameCalendarDate(a time.Time, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()

	return ay == by && am == bm && ad == bd
}
