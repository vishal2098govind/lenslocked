package models

import (
	"database/sql"
	"time"

	tokenM "github.com/vishal2098govind/lenslocked/models/token"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB

	TokenManager tokenM.TokenManager
	// ResetDuration is the amount of time that a PasswordReset is valid for.
	ResetDuration time.Duration
}
