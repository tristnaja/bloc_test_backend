package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string) (string, error) {
	tokenKey := os.Getenv("JWT_SECRET_KEY")

	if tokenKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set")
	}

	jwtKey := []byte(tokenKey)

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}
