package models

import (
	"fmt"
	"time"
)

type CreatePasswordResetRequest struct {
	UserID int
}

type CreatePasswordResetResponse struct {
	PasswordReset *PasswordReset
}

func (prs *PasswordResetService) Create(r CreatePasswordResetRequest) (*CreatePasswordResetResponse, error) {
	tokenHash, token, err := prs.TokenManager.New()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	pr := PasswordReset{
		Token:     token,
		TokenHash: tokenHash,
		UserID:    r.UserID,
		ExpiresAt: time.Now().Add(prs.defaultResetDuration()),
	}

	row := prs.DB.QueryRow(`
	INSERT INTO password_resets (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id) DO
	UPDATE
	SET token_hash = $2,
		expires_at = $3
	RETURNING id;
	`, pr.UserID, pr.TokenHash, pr.ExpiresAt)
	err = row.Scan(&pr.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &CreatePasswordResetResponse{
		PasswordReset: &pr,
	}, nil
}

func (prs *PasswordResetService) defaultResetDuration() time.Duration {
	if prs.ResetDuration == 0 {
		return DefaultResetDuration
	}
	return prs.ResetDuration
}
