package dao

import (
	"errors"
	"testing"

	"backend/domain"

	"gorm.io/gorm"
)

func TestNewUserDAO(t *testing.T) {
	db := newDAOTestDB(t)

	userDAO := NewUserDAO(db)

	if userDAO == nil {
		t.Fatal("expected user DAO")
	}
}

func TestUserDAOCreateAndFindByEmail(t *testing.T) {
	db := newDAOTestDB(t)
	userDAO := NewUserDAO(db)

	user := &domain.User{
		Name:     "Test User",
		Email:    "test@mail.com",
		Password: "hashed-password",
		Role:     domain.UserRoleClient,
	}
	if err := userDAO.Create(user); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	found, err := userDAO.FindByEmail("test@mail.com")
	if err != nil {
		t.Fatalf("FindByEmail returned error: %v", err)
	}

	if found.ID == 0 || found.Email != user.Email {
		t.Fatalf("unexpected found user: %+v", found)
	}
}

func TestUserDAOFindByEmailMissingUser(t *testing.T) {
	db := newDAOTestDB(t)
	userDAO := NewUserDAO(db)

	_, err := userDAO.FindByEmail("missing@mail.com")
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}
}

func TestUserDAOCreateDuplicateEmail(t *testing.T) {
	db := newDAOTestDB(t)
	userDAO := NewUserDAO(db)
	user := &domain.User{Name: "Test User", Email: "test@mail.com", Password: "password", Role: domain.UserRoleClient}

	if err := userDAO.Create(user); err != nil {
		t.Fatalf("first Create returned error: %v", err)
	}

	duplicate := &domain.User{Name: "Other User", Email: "test@mail.com", Password: "password", Role: domain.UserRoleClient}
	if err := userDAO.Create(duplicate); err == nil {
		t.Fatal("expected duplicate email error")
	}
}
