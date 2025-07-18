package models

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticateRequest struct {
	Email    string
	Password string
}

type AuthenticateResponse struct {
	User *User
}

var (
	ErrUserNotFound       = fmt.Errorf("user not found")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
)

func (us *UserService) Authenticate(r AuthenticateRequest) (*AuthenticateResponse, error) {
	row := us.DB.QueryRow(`
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`, r.Email)

	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &AuthenticateResponse{User: &user}, nil
}
