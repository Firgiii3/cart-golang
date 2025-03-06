package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Struct User
type User struct {
	gorm.Model
	Username   string `json:"username"`
	Password   string `json:"password"`
	Permission string `json:"permission"` // Tambahkan field permission
}

// Fungsi untuk hash password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Fungsi untuk verifikasi password
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
