package auth

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Hash password", slog.String("Error", err.Error()))
		return "", err
	}

	return string(hash), nil
}

func ValidPassword(password string, hash string) bool {
	pass := []byte(password)
	h := []byte(hash)
	err := bcrypt.CompareHashAndPassword(h, pass)

	return err == nil
}
