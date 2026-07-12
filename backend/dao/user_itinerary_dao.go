package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type UserItineraryDAO struct {
	db *gorm.DB
}

func NewUserItineraryDAO(db *gorm.DB) *UserItineraryDAO {
	return &UserItineraryDAO{db: db}
}

func (dao *UserItineraryDAO) HasActiveTicket(userID uint, eventID uint) (bool, error) {
	var count int64

	if err := dao.db.Model(&domain.Ticket{}).
		Where("user_id = ? AND event_id = ? AND status = ?", userID, eventID, domain.TicketStatusActive).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (dao *UserItineraryDAO) GetByUserAndEvent(userID uint, eventID uint) (*domain.UserItinerary, error) {
	var itinerary domain.UserItinerary

	if err := dao.db.Where("user_id = ? AND event_id = ?", userID, eventID).First(&itinerary).Error; err != nil {
		return nil, err
	}

	return &itinerary, nil
}

func (dao *UserItineraryDAO) Create(itinerary *domain.UserItinerary) error {
	return dao.db.Create(itinerary).Error
}

func (dao *UserItineraryDAO) GetItemsByItineraryID(itineraryID uint) ([]domain.ItineraryItem, error) {
	var items []domain.ItineraryItem

	if err := dao.db.Where("user_itinerary_id = ?", itineraryID).Order("start_time asc").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (dao *UserItineraryDAO) CreateItem(item *domain.ItineraryItem) error {
	return dao.db.Create(item).Error
}

func (dao *UserItineraryDAO) GetItemByID(id uint) (*domain.ItineraryItem, error) {
	var item domain.ItineraryItem

	if err := dao.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (dao *UserItineraryDAO) DeleteItem(id uint) error {
	return dao.db.Delete(&domain.ItineraryItem{}, id).Error
}

func (dao *UserItineraryDAO) DeleteByUserAndEvent(userID uint, eventID uint) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		var itineraryIDs []uint
		if err := tx.Model(&domain.UserItinerary{}).
			Where("user_id = ? AND event_id = ?", userID, eventID).
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
	})
}
