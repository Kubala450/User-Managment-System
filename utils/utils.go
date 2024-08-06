package utils

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func ValidateUserInput(username, password string) error {
	trimmedUsername := strings.TrimSpace(username)
	trimmedPassword := strings.TrimSpace(password)

	if trimmedUsername == "" {
		return errors.New("username cannot be empty or consist only of whitespace")
	}

	if trimmedPassword == "" {
		return errors.New("password cannot be empty or consist only of whitespace")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}