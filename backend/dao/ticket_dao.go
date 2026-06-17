package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type TicketDAO struct {
	db *gorm.DB
}

func NewTicketDAO(db *gorm.DB) *TicketDAO {
	return &TicketDAO{db: db}
}

func (dao *TicketDAO) GetEventByID(id uint) (*domain.Event, error) {
	var event domain.Event

	if err := dao.db.First(&event, id).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (dao *TicketDAO) CountTickets() (int64, error) {
	var count int64

	if err := dao.db.Model(&domain.Ticket{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *TicketDAO) GetByUserID(userID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket

	if err := dao.db.
		Preload("Event").
		Where("user_id = ?", userID).
		Order("purchase_date desc").
		Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func (dao *TicketDAO) GetByID(id uint) (*domain.Ticket, error) {
	var ticket domain.Ticket

	if err := dao.db.Preload("Event").First(&ticket, id).Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (dao *TicketDAO) CreatePurchase(ticket *domain.Ticket, event *domain.Event) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(event).Error; err != nil {
			return err
		}

		if err := tx.Create(ticket).Error; err != nil {
			return err
		}

		return nil
	})
}

func (dao *TicketDAO) CancelTicket(ticket *domain.Ticket, event *domain.Event) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(event).Error; err != nil {
			return err
		}

		if err := tx.Save(ticket).Error; err != nil {
			return err
		}

		return nil
	})
}
