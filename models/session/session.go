package models

import "database/sql"

const (
	// The min number of bytes to be used for each session token
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int

	// only set while creating new session.
	// When looking up a session, this will
	// be left empty, as we only store hash
	// of a session token in the DB and we
	// cannot reverse it into a raw token
	Token string

	TokenHash string
}

type TokenManager struct {
	// BytesPerToken is used to determine how many bytes to use while
	// generating each session token. If this value is not set or is
	// less than MinBytesPerToken const it will be ignored and
	// MinBytesPerToken will be used.
	BytesPerToken int
}

type SessionService struct {
	DB *sql.DB

	TokenManager TokenManager
}
