package model

import (
	"taskboard/common"

	"golang.org/x/crypto/bcrypt"
)

// User is user of the app.
type User struct {
	ID           string `gorm:"primary_key;size:32"`
	Name         string `gorm:"size:255;not null;unique"`
	PasswordHash string `gorm:"size:255;not null;"`
	Avator       string `gorm:"size:255"`
	Version      int    `gorm:"not null"` // Version for optimistic lock
}

// NewUser returns created new user
func NewUser(name, rawpassword, avator string) *User {
	result := &User{
		ID:      "user_" + common.GenerateID(),
		Name:    name,
		Avator:  avator,
		Version: 1,
	}
	result.SetPassword(rawpassword)
	return result
}

// SetPassword sets the hash of specified password to user
func (user *User) SetPassword(password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.PasswordHash = string(hash)
}

// VerifyPassword checks whether specified password matches PasswordHash in database.
func (user *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}
