package models

import (
	"fmt"
)

type CreateSessionRequest struct {
	UserID int
}

type CreateSessionResponse struct {
	Session *Session
}

func (ss *SessionService) Create(r CreateSessionRequest) (*CreateSessionResponse, error) {

	tokenHash, token, err := ss.TokenManager.New()
	if err != nil {
		return nil, err
	}
	session := Session{
		UserID:    r.UserID,
		TokenHash: tokenHash,
		Token:     token,
	}

	row := ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2)
		ON CONFLICT(user_id) 
		DO UPDATE SET token_hash=$2
		RETURNING id
	`, r.UserID, tokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &CreateSessionResponse{
		Session: &session,
	}, nil

}
