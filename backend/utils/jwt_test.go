package utils

import "testing"

func TestGenerateJWTGeneratesValidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES", "1")

	token, err := GenerateJWT(1, "Test User", "test@mail.com")
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
}

func TestValidateJWTRejectsInvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	if _, err := ValidateJWT("invalid-token"); err == nil {
		t.Fatal("expected invalid token error")
	}
}
