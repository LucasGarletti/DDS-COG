package dao

import (
	"backend/domain"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	if err := dao.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (dao *UserDAO) Create(user *domain.User) error {
	return dao.db.Create(user).Error
}
