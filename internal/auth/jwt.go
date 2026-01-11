package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/arsenh/DevLore/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUser struct {
	ID       int
	Email    string
	FullName string
}

func GenerateJWTToken(userID int, email string, fullName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        userID,
		"email":     email,
		"full_name": fullName,
		"exp":       time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(config.JWTSecretKey)

	return tokenString, err
}

func ParseJWTToken(tokenString string) (*AuthUser, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// Prevent "alg:none" attack
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return config.JWTSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Extract fields
	id, ok := claims["id"].(float64) // JSON numbers are float64
	if !ok {
		return nil, errors.New("id missing")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email missing")
	}

	fullName, ok := claims["full_name"].(string)
	if !ok {
		return nil, errors.New("full_name missing")
	}

	return &AuthUser{
		ID:       int(id),
		FullName: fullName,
		Email:    email,
	}, nil
}
