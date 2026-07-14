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
	createErr    error
	createdUser  *domain.User
}

func (repo *fakeRegisterUserRepository) FindByEmail(email string) (*domain.User, error) {
	if repo.findErr != nil {
		return nil, repo.findErr
	}

	return repo.existingUser, nil
}

func (repo *fakeRegisterUserRepository) Create(user *domain.User) error {
	if repo.createErr != nil {
		return repo.createErr
	}
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
			Role:     domain.UserRoleClient,
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

func TestRegisterCreatesClientUser(t *testing.T) {
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

	if user.Role != domain.UserRoleClient {
		t.Fatalf("expected registered user role client, got %s", user.Role)
	}

	if repo.createdUser.Role != domain.UserRoleClient {
		t.Fatalf("expected created user role client, got %s", repo.createdUser.Role)
	}
}

func TestRegisterDoesNotAllowAdminRole(t *testing.T) {
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

	if user.Role == domain.UserRoleAdmin {
		t.Fatal("expected public register not to create admin user")
	}
}

func TestLoginWithValidCredentialsReturnsToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES", "1")

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
			Role:     domain.UserRoleClient,
		},
	})

	login, err := service.Login(LoginInput{
		Email:    "test@mail.com",
		Password: "correct-password",
	})
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if login.Token == "" {
		t.Fatal("expected token")
	}

	if login.User.Role != domain.UserRoleClient {
		t.Fatalf("expected user role client, got %s", login.User.Role)
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

func TestRegisterWithFindErrorReturnsError(t *testing.T) {
	findErr := errors.New("find error")
	service := NewAuthService(&fakeRegisterUserRepository{findErr: findErr})

	_, err := service.Register(RegisterInput{
		Name:     "Test User",
		Email:    "test@mail.com",
		Password: "password",
	})
	if !errors.Is(err, findErr) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestRegisterWithCreateErrorReturnsError(t *testing.T) {
	createErr := errors.New("create error")
	service := NewAuthService(&fakeRegisterUserRepository{
		findErr:   gorm.ErrRecordNotFound,
		createErr: createErr,
	})

	_, err := service.Register(RegisterInput{
		Name:     "Test User",
		Email:    "test@mail.com",
		Password: "password",
	})
	if !errors.Is(err, createErr) {
		t.Fatalf("expected create error, got %v", err)
	}
}

func TestLoginWithRepositoryErrorReturnsError(t *testing.T) {
	repoErr := errors.New("repo error")
	service := NewAuthService(fakeUserRepository{err: repoErr})

	_, err := service.Login(LoginInput{
		Email:    "test@mail.com",
		Password: "password",
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestLoginWithTokenGenerationErrorReturnsError(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_EXPIRES", "1")

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
			Role:     domain.UserRoleClient,
		},
	})

	if _, err := service.Login(LoginInput{Email: "test@mail.com", Password: "correct-password"}); err == nil {
		t.Fatal("expected token generation error")
	}
}
