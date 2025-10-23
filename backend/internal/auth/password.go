package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// ErrInvalidPassword indicates that the supplied password does not match the stored hash.
var ErrInvalidPassword = errors.New("invalid password")

// HashPassword returns a bcrypt hash for the supplied plaintext password.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword compares a bcrypt hash against the supplied plaintext password.
func ComparePassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidPassword
	}
	return err
}
