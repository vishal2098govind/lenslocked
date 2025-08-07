package models

import (
	"fmt"
)

type deletePasswordResetTokenRequest struct {
	Token string
}

func (prs *PasswordResetService) delete(r deletePasswordResetTokenRequest) error {
	tokenHash := prs.TokenManager.Hash(r.Token)
	_, err := prs.DB.Exec(`
		DELETE FROM password_resets
		WHERE token_hash = $1`, tokenHash)

	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
