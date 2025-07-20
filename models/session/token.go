package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/vishal2098govind/lenslocked/rand"
)

func (tm TokenManager) New() (token string, tokenHash string, err error) {
	bytesPerToken := max(tm.BytesPerToken, MinBytesPerToken)
	token, err = rand.String(bytesPerToken)
	if err != nil {
		log.Println(err)
		return "", "", fmt.Errorf("create: %w", err)
	}

	tokenHash = tm.hash(token)
	return tokenHash, token, nil
}

func (tm *TokenManager) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
