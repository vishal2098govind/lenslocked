package models

import (
	"database/sql"
	"errors"
	"fmt"
)

func (us *UserService) ViaEmail(email string) (*User, error) {
	row := us.DB.QueryRow(`
	SELECT id, email, password_hash
	FROM users
	WHERE email = $1
	`, email)

	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user via email: %w", err)
	}

	return &user, nil
}
