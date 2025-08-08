package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Email    string
	Password string
}

type CreateUserResponse struct {
	User *User
}

func (us *UserService) Create(r CreateUserRequest) (*CreateUserResponse, error) {

	email := strings.ToLower(r.Email)

	passh, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	hashedP := string(passh)

	row := us.DB.QueryRow(`
	INSERT INTO users (email, password_hash)
	VALUES ($1, $2)
	RETURNING id;
	`, email, hashedP)

	user := User{
		Email:        email,
		PasswordHash: hashedP,
	}

	err = row.Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailAlreadyExists
			}
		}
		return nil, fmt.Errorf("create user:  %w", err)
	}

	return &CreateUserResponse{
		User: &user,
	}, nil
}
