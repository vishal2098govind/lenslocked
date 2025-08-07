package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/vishal2098govind/lenslocked/rand"
)

const (
	// The min number of bytes to be used for each session token
	MinBytesPerToken = 32
)

type TokenManager struct {
	// BytesPerToken is used to determine how many bytes to use while
	// generating each session token. If this value is not set or is
	// less than MinBytesPerToken const it will be ignored and
	// MinBytesPerToken will be used.
	BytesPerToken int
}

func (tm TokenManager) New() (tokenHash string, token string, err error) {
	bytesPerToken := max(tm.BytesPerToken, MinBytesPerToken)
	token, err = rand.String(bytesPerToken)
	if err != nil {
		log.Println(err)
		return "", "", fmt.Errorf("create: %w", err)
	}

	tokenHash = tm.Hash(token)
	return tokenHash, token, nil
}

func (tm *TokenManager) Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
