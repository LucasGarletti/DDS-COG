package utils

import (
	"testing"
	"time"

	"backend/domain"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWTGeneratesValidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES", "1")

	token, err := GenerateJWT(1, "Test User", "test@mail.com", domain.UserRoleClient)
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %v", err)
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if claims.ID != 1 {
		t.Fatalf("expected user id 1, got %d", claims.ID)
	}

	if claims.Name != "Test User" {
		t.Fatalf("expected name Test User, got %s", claims.Name)
	}

	if claims.Email != "test@mail.com" {
		t.Fatalf("expected email test@mail.com, got %s", claims.Email)
	}

	if claims.Role != domain.UserRoleClient {
		t.Fatalf("expected role client, got %s", claims.Role)
	}
}

func TestGenerateJWTRejectsInvalidRole(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES", "1")

	if _, err := GenerateJWT(1, "Test User", "test@mail.com", "superadmin"); err == nil {
		t.Fatal("expected invalid role error")
	}
}

func TestValidateJWTRejectsInvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	if _, err := ValidateJWT("invalid-token"); err == nil {
		t.Fatal("expected invalid token error")
	}
}

func TestValidateJWTRejectsInvalidRole(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	claims := JWTClaims{
		ID:    1,
		Name:  "Test User",
		Email: "test@mail.com",
		Role:  "superadmin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("could not sign token: %v", err)
	}

	if _, err := ValidateJWT(token); err == nil {
		t.Fatal("expected invalid role error")
	}
}

func TestValidateJWTRejectsEmptyRole(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	claims := JWTClaims{
		ID:    1,
		Name:  "Test User",
		Email: "test@mail.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("could not sign token: %v", err)
	}

	if _, err := ValidateJWT(token); err == nil {
		t.Fatal("expected empty role error")
	}
}

func TestValidateJWTRejectsMissingSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	if _, err := ValidateJWT("token"); err == nil {
		t.Fatal("expected missing secret error")
	}
}

func TestValidateJWTRejectsEmptyToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	if _, err := ValidateJWT(" "); err == nil {
		t.Fatal("expected empty token error")
	}
}

func TestIsValidUserRole(t *testing.T) {
	tests := []struct {
		name  string
		role  string
		valid bool
	}{
		{name: "client", role: domain.UserRoleClient, valid: true},
		{name: "admin", role: domain.UserRoleAdmin, valid: true},
		{name: "empty", role: "", valid: false},
		{name: "unknown", role: "superadmin", valid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if IsValidUserRole(test.role) != test.valid {
				t.Fatalf("expected role %q valid=%t", test.role, test.valid)
			}
		})
	}
}
