package models

import (
	"database/sql"
	"fmt"
	"log"

	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type GetUserIdRequest struct {
	Token string
}

type GetUserIdResponse struct {
	User *userM.User
}

func (ss *SessionService) User(r GetUserIdRequest) (*GetUserIdResponse, error) {
	tokenHash := ss.TokenManager.Hash(r.Token)

	user := userM.User{}
	row := ss.DB.QueryRow(`
		SELECT u.id, u.email
		FROM sessions s
		JOIN users u on s.user_id = u.id
		WHERE token_hash = $1
	`, tokenHash)

	err := row.Scan(
		&user.ID,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		// invalid session
		return &GetUserIdResponse{
			User: nil,
		}, nil
	}

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("user: %w", err)
	}

	return &GetUserIdResponse{
		User: &user,
	}, nil
}
