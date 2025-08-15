package models

import "fmt"

type DeleteGalleryRequest struct {
	ID int
}

func (gs *GalleryService) DeleteByID(req DeleteGalleryRequest) error {

	_, err := gs.DB.Exec(`DELETE FROM galleries WHERE id = $1`, req.ID)
	if err != nil {
		return fmt.Errorf("delete gallery by id: %w", err)
	}

	return nil
}
