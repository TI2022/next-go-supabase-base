package entity

import "time"

type UserID string

type User struct {
	ID           UserID
	Email        string
	PasswordHash string
	Name         *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

