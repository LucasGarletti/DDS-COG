package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type EventDAO struct {
	db *gorm.DB
}

func NewEventDAO(db *gorm.DB) *EventDAO {
	return &EventDAO{db: db}
}

func (dao *EventDAO) GetAll() ([]domain.Event, error) {
	var events []domain.Event

	if err := dao.db.Where("status = ?", domain.EventStatusActive).Order("date asc").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
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
