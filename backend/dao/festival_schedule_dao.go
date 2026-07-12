package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type FestivalScheduleDAO struct {
	db *gorm.DB
}

func NewFestivalScheduleDAO(db *gorm.DB) *FestivalScheduleDAO {
	return &FestivalScheduleDAO{db: db}
}

func (dao *FestivalScheduleDAO) Create(schedule *domain.FestivalSchedule) error {
	return dao.db.Create(schedule).Error
}

func (dao *FestivalScheduleDAO) GetByID(id uint) (*domain.FestivalSchedule, error) {
	var schedule domain.FestivalSchedule

	if err := dao.db.First(&schedule, id).Error; err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (dao *FestivalScheduleDAO) GetByEventID(eventID uint) ([]domain.FestivalSchedule, error) {
	var schedules []domain.FestivalSchedule

	if err := dao.db.Where("event_id = ?", eventID).Order("start_time asc").Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (dao *FestivalScheduleDAO) Delete(id uint) error {
	return dao.db.Delete(&domain.FestivalSchedule{}, id).Error
}
