package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordRequest struct {
	UserID      int
	NewPassword string
}

type UpdatePasswordResponse struct {
	User *User
}

func (us *UserService) UpdatePassword(r UpdatePasswordRequest) (*UpdatePasswordResponse, error) {
	passh, err := bcrypt.GenerateFromPassword([]byte(r.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}

	row := us.DB.QueryRow(`
		UPDATE users SET
		password_hash = $1
		WHERE id = $2
		RETURNING id, email, password_hash;
	`, passh, r.UserID)

	user := User{}
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("update password: %w", err)
	}

	return &UpdatePasswordResponse{
		User: &user,
	}, nil
}
