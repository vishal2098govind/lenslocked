package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type GetGalleriesByUserIDRequest struct {
	UserID int
}

type GetGalleriesResponse struct {
	Galleries []Gallery
}

func (gs *GalleryService) ByUserID(req GetGalleriesByUserIDRequest) (*GetGalleriesResponse, error) {
	rows, err := gs.DB.Query(`
	 SELECT id, user_id, title
	 FROM galleries
	 WHERE user_id = $1
	`, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("get galleries by user id: %w", err)
	}

	galleries := []Gallery{}
	for rows.Next() {
		var g Gallery
		rows.Scan(&g.ID, &g.UserID, &g.Title)
		galleries = append(galleries, g)
	}

	return &GetGalleriesResponse{
		Galleries: galleries,
	}, nil
}

type GetGalleryByIDRequest struct {
	ID int
}

type GetGalleryResponse struct {
	Gallery *Gallery
}

func (gs *GalleryService) ByID(req GetGalleryByIDRequest) (*GetGalleryResponse, error) {
	row := gs.DB.QueryRow(`
	 SELECT id, user_id, title
	 FROM galleries
	 WHERE id = $1
	`, req.ID)

	var g Gallery
	err := row.Scan(&g.ID, &g.UserID, &g.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGalleryNotFound
		}
		return nil, fmt.Errorf("by user id: %w", err)
	}

	return &GetGalleryResponse{
		Gallery: &g,
	}, nil
}
