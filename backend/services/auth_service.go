package services

import (
	"errors"

	"backend/dao"
	"backend/domain"
	"backend/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
	User  *domain.User
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

func (service *AuthService) Login(input LoginInput) (*LoginOutput, error) {
	user, err := service.userDAO.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token: token,
		User:  user,
	}, nil
}
