package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

var (
	ErrUserNotFound       = fmt.Errorf("models/user: user not found")
	ErrInvalidCredentials = fmt.Errorf("models/user: invalid credentials")
	ErrEmailAlreadyExists = fmt.Errorf("models/user: email already exists")
)
