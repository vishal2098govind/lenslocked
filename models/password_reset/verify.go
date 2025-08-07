package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type VerifyRequest struct {
	Token string
}

type VerifyResponse struct {
	PasswordReset *PasswordReset
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

func (prs *PasswordResetService) Verify(r VerifyRequest) (*VerifyResponse, error) {
	tokenHash := prs.TokenManager.Hash(r.Token)
	row := prs.DB.QueryRow(`
	SELECT id, user_id, token_hash, expires_at
	FROM password_resets
	WHERE token_hash = $1`, tokenHash)

	pr := PasswordReset{}
	err := row.Scan(&pr.ID, &pr.UserID, &pr.TokenHash, &pr.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}
	if pr.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}

	err = prs.delete(deletePasswordResetTokenRequest(r))
	if err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}

	return &VerifyResponse{
		PasswordReset: &pr,
	}, nil
}
