package models

import "fmt"

type DeleteSessionRequest struct {
	Token string
}

type DeleteSessionResponse struct {
}

func (ss *SessionService) Delete(r *DeleteSessionRequest) (*DeleteSessionResponse, error) {
	tokenHash := ss.TokenManager.hash(r.Token)
	_, err := ss.DB.Exec(`
		DELETE FROM sessions 
		WHERE token_hash = $1`, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("delete: %w", err)
	}

	return &DeleteSessionResponse{}, nil
}
