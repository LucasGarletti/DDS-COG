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

	if err := dao.db.Order("date asc").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (dao *EventDAO) GetByID(id uint) (*domain.Event, error) {
	var event domain.Event

	if err := dao.db.First(&event, id).Error; err != nil {
		return nil, err
	}

	return &event, nil
}
