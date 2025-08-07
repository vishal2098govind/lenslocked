package models

import (
	"database/sql"

	tokenM "github.com/vishal2098govind/lenslocked/models/token"
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

type SessionService struct {
	DB *sql.DB

	TokenManager tokenM.TokenManager
}
