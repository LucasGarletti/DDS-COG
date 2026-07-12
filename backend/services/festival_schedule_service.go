package services

import (
	"errors"
	"strings"
	"time"

	"backend/domain"
)

var (
	ErrInvalidFestivalSchedule = errors.New("invalid festival schedule")
	ErrEventNotFestival        = errors.New("event is not festival")
)

type FestivalEventRepository interface {
	GetByIDForAdmin(id uint) (*domain.Event, error)
}

type FestivalScheduleRepository interface {
	Create(schedule *domain.FestivalSchedule) error
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
