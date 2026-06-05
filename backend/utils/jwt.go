package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, name string, email string) (string, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return "", errors.New("JWT_SECRET is required")
	}

	expiresInHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRES"))
	if err != nil || expiresInHours <= 0 {
		return "", errors.New("JWT_EXPIRES must be a positive number of hours")
	}

	now := time.Now()
	claims := JWTClaims{
		ID:    userID,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiresInHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
