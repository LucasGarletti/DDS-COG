package services

import (
	"errors"
	"testing"

	"backend/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type fakeUserRepository struct {
	user *domain.User
	err  error
}

func (repo fakeUserRepository) FindByEmail(email string) (*domain.User, error) {
	if repo.err != nil {
		return nil, repo.err
	}

	return repo.user, nil
}

func (repo fakeUserRepository) Create(user *domain.User) error {
	return nil
}

func TestLoginWithInvalidCredentialsReturnsError(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("could not hash password: %v", err)
	}

	service := NewAuthService(fakeUserRepository{
		user: &domain.User{
			ID:       1,
			Name:     "Test User",
			Email:    "test@mail.com",
			Password: string(hashedPassword),
		},
	})

	_, err = service.Login(LoginInput{
		Email:    "test@mail.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginWithMissingUserReturnsError(t *testing.T) {
	service := NewAuthService(fakeUserRepository{err: gorm.ErrRecordNotFound})

	_, err := service.Login(LoginInput{
		Email:    "missing@mail.com",
		Password: "password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
