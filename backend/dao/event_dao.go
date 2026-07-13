package dao

import (
	"backend/domain"
	"strings"

	"gorm.io/gorm"
)

type EventDAO struct {
	db *gorm.DB
}

func NewEventDAO(db *gorm.DB) *EventDAO {
	return &EventDAO{db: db}
}

func (dao *EventDAO) GetAll(filters domain.EventFilters) ([]domain.Event, error) {
	var events []domain.Event

	query := dao.db.Where("status = ?", domain.EventStatusActive)

	if filters.Search != "" {
		search := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("(LOWER(title) LIKE ? OR LOWER(description) LIKE ?)", search, search)
	}

	if filters.Location != "" {
		query = query.Where("LOWER(location) LIKE ?", "%"+strings.ToLower(filters.Location)+"%")
	}

	if filters.DateFrom != nil {
		query = query.Where("date >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("date <= ?", *filters.DateTo)
	}

	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}

	if filters.IsFestival != nil {
		query = query.Where("is_festival = ?", *filters.IsFestival)
	}

	if filters.AvailableOnly {
		query = query.Where("available_capacity > ?", 0)
	}

	query = query.Order(eventSortOrder(filters.Sort))

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func eventSortOrder(sort string) string {
	switch sort {
	case domain.EventSortDateDesc:
		return "date desc"
	case domain.EventSortPriceAsc:
		return "price asc"
	case domain.EventSortPriceDesc:
		return "price desc"
	default:
		return "date asc"
	}
}

func (dao *EventDAO) GetByID(id uint) (*domain.Event, error) {
	var event domain.Event

	if err := dao.db.Where("status = ?", domain.EventStatusActive).First(&event, id).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (dao *EventDAO) GetByIDForAdmin(id uint) (*domain.Event, error) {
	var event domain.Event

	if err := dao.db.First(&event, id).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (dao *EventDAO) Create(event *domain.Event) error {
	return dao.db.Create(event).Error
}

func (dao *EventDAO) Save(event *domain.Event) error {
	return dao.db.Save(event).Error
}

func (dao *EventDAO) CancelWithCleanup(id uint) (*domain.Event, error) {
	var event domain.Event

	err := dao.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&event, id).Error; err != nil {
			return err
		}

		if event.IsFestival {
			if err := deleteFestivalItineraries(tx, id); err != nil {
				return err
			}

			if err := tx.Where("event_id = ?", id).Delete(&domain.FestivalSchedule{}).Error; err != nil {
				return err
			}
		}

		event.Status = domain.EventStatusCancelled
		return tx.Save(&event).Error
	})
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func deleteFestivalItineraries(tx *gorm.DB, eventID uint) error {
	var itineraryIDs []uint
	if err := tx.Model(&domain.UserItinerary{}).
		Where("event_id = ?", eventID).
		Pluck("id", &itineraryIDs).Error; err != nil {
		return err
	}

	if len(itineraryIDs) == 0 {
		return nil
	}

	if err := tx.Where("user_itinerary_id IN ?", itineraryIDs).Delete(&domain.ItineraryItem{}).Error; err != nil {
		return err
	}

	return tx.Where("id IN ?", itineraryIDs).Delete(&domain.UserItinerary{}).Error
}
