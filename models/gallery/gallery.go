package models

import (
	"database/sql"
	"fmt"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

func (g *Gallery) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"title": g.Title,
		"id":    g.ID,
	}
}

type GalleryService struct {
	DB *sql.DB
}

var (
	ErrGalleryNotFound = fmt.Errorf("gallery not found")
)
