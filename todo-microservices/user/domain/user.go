package user

import (
	"time"
)

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	Password  []byte
	CreatedAt *time.Time
}

func NewUser(email string, password []byte) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}
