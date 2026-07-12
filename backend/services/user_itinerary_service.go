package services

import (
	"errors"
	"strings"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

var (
	ErrItineraryForbidden    = errors.New("user does not have active ticket for event")
	ErrInvalidItineraryItem  = errors.New("invalid itinerary item")
	ErrItineraryItemNotOwned = errors.New("itinerary item does not belong to user itinerary")
	ErrScheduleEventMismatch = errors.New("schedule does not belong to event")
)

type UserItineraryRepository interface {
	HasActiveTicket(userID uint, eventID uint) (bool, error)
	GetByUserAndEvent(userID uint, eventID uint) (*domain.UserItinerary, error)
	Create(itinerary *domain.UserItinerary) error
	GetItemsByItineraryID(itineraryID uint) ([]domain.ItineraryItem, error)
	CreateItem(item *domain.ItineraryItem) error
	GetItemByID(id uint) (*domain.ItineraryItem, error)
	DeleteItem(id uint) error
}

type UserItineraryService struct {
	eventDAO     FestivalEventRepository
	scheduleDAO  FestivalScheduleRepository
	itineraryDAO UserItineraryRepository
}

type AddPersonalActivityInput struct {
	UserID    uint
	EventID   uint
	Title     string
	Location  string
	Notes     string
	StartTime time.Time
	EndTime   time.Time
}

type ItineraryResult struct {
	Event        *domain.Event          `json:"event"`
	Items        []domain.ItineraryItem `json:"items"`
	HasConflicts bool                   `json:"has_conflicts"`
}

type AddItineraryItemResult struct {
	Item        *domain.ItineraryItem `json:"item"`
	HasConflict bool                  `json:"has_conflict"`
}

func NewUserItineraryService(eventDAO FestivalEventRepository, scheduleDAO FestivalScheduleRepository, itineraryDAO UserItineraryRepository) *UserItineraryService {
	return &UserItineraryService{
		eventDAO:     eventDAO,
		scheduleDAO:  scheduleDAO,
		itineraryDAO: itineraryDAO,
	}
}

func (service *UserItineraryService) GetItinerary(userID uint, eventID uint) (*ItineraryResult, error) {
	event, itinerary, err := service.requireFestivalAccess(userID, eventID)
	if err != nil {
		return nil, err
	}

	items, err := service.itineraryDAO.GetItemsByItineraryID(itinerary.ID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []domain.ItineraryItem{}
	}

	return &ItineraryResult{
		Event:        event,
		Items:        items,
		HasConflicts: hasAnyTimeConflict(items),
	}, nil
}

func (service *UserItineraryService) AddShow(userID uint, eventID uint, scheduleID uint) (*AddItineraryItemResult, error) {
	if scheduleID == 0 {
		return nil, ErrInvalidItineraryItem
	}

	_, itinerary, err := service.requireFestivalAccess(userID, eventID)
	if err != nil {
		return nil, err
	}

	schedule, err := service.scheduleDAO.GetByID(scheduleID)
	if err != nil {
		return nil, err
	}

	if schedule.EventID != eventID {
		return nil, ErrScheduleEventMismatch
	}

	item := &domain.ItineraryItem{
		UserItineraryID:    itinerary.ID,
		FestivalScheduleID: &schedule.ID,
		Type:               domain.ItineraryItemTypeShow,
		Title:              schedule.Artist,
		Location:           schedule.Stage,
		StartTime:          schedule.StartTime,
		EndTime:            schedule.EndTime,
	}

	return service.createItemWithConflictStatus(itinerary.ID, item)
}

func (service *UserItineraryService) AddPersonalActivity(input AddPersonalActivityInput) (*AddItineraryItemResult, error) {
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" || input.StartTime.IsZero() || input.EndTime.IsZero() || !input.EndTime.After(input.StartTime) {
		return nil, ErrInvalidItineraryItem
	}

	_, itinerary, err := service.requireFestivalAccess(input.UserID, input.EventID)
	if err != nil {
		return nil, err
	}

	item := &domain.ItineraryItem{
		UserItineraryID: itinerary.ID,
		Type:            domain.ItineraryItemTypePersonal,
		Title:           input.Title,
		Location:        strings.TrimSpace(input.Location),
		Notes:           input.Notes,
		StartTime:       input.StartTime,
		EndTime:         input.EndTime,
	}

	return service.createItemWithConflictStatus(itinerary.ID, item)
}

func (service *UserItineraryService) DeleteItem(userID uint, eventID uint, itemID uint) error {
	if itemID == 0 {
		return ErrInvalidItineraryItem
	}

	_, itinerary, err := service.requireFestivalAccess(userID, eventID)
	if err != nil {
		return err
	}

	item, err := service.itineraryDAO.GetItemByID(itemID)
	if err != nil {
		return err
	}

	if item.UserItineraryID != itinerary.ID {
		return ErrItineraryItemNotOwned
	}

	return service.itineraryDAO.DeleteItem(itemID)
}

func (service *UserItineraryService) createItemWithConflictStatus(itineraryID uint, item *domain.ItineraryItem) (*AddItineraryItemResult, error) {
	items, err := service.itineraryDAO.GetItemsByItineraryID(itineraryID)
	if err != nil {
		return nil, err
	}

	hasConflict := hasTimeConflict(*item, items)

	if err := service.itineraryDAO.CreateItem(item); err != nil {
		return nil, err
	}

	return &AddItineraryItemResult{
		Item:        item,
		HasConflict: hasConflict,
	}, nil
}

func (service *UserItineraryService) requireFestivalAccess(userID uint, eventID uint) (*domain.Event, *domain.UserItinerary, error) {
	if userID == 0 || eventID == 0 {
		return nil, nil, ErrInvalidItineraryItem
	}

	event, err := service.eventDAO.GetByIDForAdmin(eventID)
	if err != nil {
		return nil, nil, err
	}

	if !event.IsFestival {
		return nil, nil, ErrEventNotFestival
	}

	hasActiveTicket, err := service.itineraryDAO.HasActiveTicket(userID, eventID)
	if err != nil {
		return nil, nil, err
	}

	if !hasActiveTicket {
		return nil, nil, ErrItineraryForbidden
	}

	itinerary, err := service.itineraryDAO.GetByUserAndEvent(userID, eventID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}

		itinerary = &domain.UserItinerary{
			UserID:  userID,
			EventID: eventID,
		}
		if err := service.itineraryDAO.Create(itinerary); err != nil {
			return nil, nil, err
		}
	}

	return event, itinerary, nil
}

func hasAnyTimeConflict(items []domain.ItineraryItem) bool {
	for i := range items {
		for j := i + 1; j < len(items); j++ {
			if overlaps(items[i].StartTime, items[i].EndTime, items[j].StartTime, items[j].EndTime) {
				return true
			}
		}
	}
	return false
}

func hasTimeConflict(candidate domain.ItineraryItem, items []domain.ItineraryItem) bool {
	for _, item := range items {
		if overlaps(candidate.StartTime, candidate.EndTime, item.StartTime, item.EndTime) {
			return true
		}
	}
	return false
}

func overlaps(aStart time.Time, aEnd time.Time, bStart time.Time, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}
