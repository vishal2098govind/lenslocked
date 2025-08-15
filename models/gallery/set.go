package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type SetGalleryTitleRequest struct {
	ID    int
	Title string
}

type SetGalleryResponse struct {
	Gallery *Gallery
}

func (gs *GalleryService) SetTitle(req SetGalleryTitleRequest) (*SetGalleryResponse, error) {
	row := gs.DB.QueryRow(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1
		RETURNING id, user_id, title
	`, req.ID, req.Title)

	var g Gallery
	err := row.Scan(&g.ID, &g.UserID, &g.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGalleryNotFound
		}

		return nil, fmt.Errorf("set gallery title: %w", err)
	}

	return &SetGalleryResponse{Gallery: &g}, nil
}
