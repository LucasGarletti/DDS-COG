package services

import (
	"errors"

	"backend/dao"
	"backend/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type AuthService struct {
	userDAO *dao.UserDAO
}

func NewAuthService(userDAO *dao.UserDAO) *AuthService {
	return &AuthService{userDAO: userDAO}
}

func (service *AuthService) Register(input RegisterInput) (*domain.User, error) {
	existingUser, err := service.userDAO.FindByEmail(input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := service.userDAO.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
