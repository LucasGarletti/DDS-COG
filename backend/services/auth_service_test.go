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

type fakeRegisterUserRepository struct {
	existingUser *domain.User
	findErr      error
	createdUser  *domain.User
}

func (repo *fakeRegisterUserRepository) FindByEmail(email string) (*domain.User, error) {
	if repo.findErr != nil {
		return nil, repo.findErr
	}

	return repo.existingUser, nil
}

func (repo *fakeRegisterUserRepository) Create(user *domain.User) error {
	repo.createdUser = user
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

func TestRegisterHashesPassword(t *testing.T) {
	repo := &fakeRegisterUserRepository{findErr: gorm.ErrRecordNotFound}
	service := NewAuthService(repo)

	user, err := service.Register(RegisterInput{
		Name:     "Test User",
		Email:    "test@mail.com",
		Password: "plain-password",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if repo.createdUser == nil {
		t.Fatal("expected user to be created")
	}

	if user.Password == "plain-password" {
		t.Fatal("expected password to be hashed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("plain-password")); err != nil {
		t.Fatalf("expected hashed password to match original password: %v", err)
	}
}

func TestRegisterWithRepeatedEmailReturnsError(t *testing.T) {
	service := NewAuthService(fakeUserRepository{
		user: &domain.User{
			ID:    1,
			Email: "test@mail.com",
		},
	})

	_, err := service.Register(RegisterInput{
		Name:     "Test User",
		Email:    "test@mail.com",
		Password: "password",
	})
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}
