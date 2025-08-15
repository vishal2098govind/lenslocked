package models

import "fmt"

type CreateGalleryRequest struct {
	UserID int
	Title  string
}

type CreateGalleryResponse struct {
	Gallery *Gallery
}

func (gs *GalleryService) Create(req CreateGalleryRequest) (*CreateGalleryResponse, error) {
	row := gs.DB.QueryRow(`
		INSERT INTO galleries (user_id, title)
		VALUES ($1, $2)
		RETURNING id
	`, req.UserID, req.Title)

	var g Gallery
	g.UserID = req.UserID
	g.Title = req.Title
	err := row.Scan(&g.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &CreateGalleryResponse{
		Gallery: &g,
	}, nil
}
