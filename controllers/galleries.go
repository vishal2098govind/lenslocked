package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vishal2098govind/lenslocked/context"
	galleryM "github.com/vishal2098govind/lenslocked/models/gallery"
)

type Galleries struct {
	GalleryService *galleryM.GalleryService
}

func (c *Galleries) New(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	type input struct {
		Title string `json:"title"`
	}
	var in input
	if err := ReadJSON(r, &in); err != nil {
		fmt.Println(err)
		WriteErrJSON(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}

	res, err := c.GalleryService.Create(galleryM.CreateGalleryRequest{
		UserID: user.ID,
		Title:  in.Title,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create gallery"})
		return
	}

	WriteJSON(w, http.StatusOK, res.Gallery.ToJSON())
}

func (c *Galleries) GetGalleriesByUserID(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	res, err := c.GalleryService.ByUserID(galleryM.GetGalleriesByUserIDRequest{
		UserID: user.ID,
	})
	if err != nil {
		WriteErrJSON(w, http.StatusInternalServerError, "Failed to get galleries")
		return
	}

	var galleries []map[string]interface{}
	for _, v := range res.Galleries {
		galleries = append(galleries, v.ToJSON())
	}

	WriteJSON(w, http.StatusOK, galleries)
}

func (c *Galleries) GetGalleryByID(w http.ResponseWriter, r *http.Request) {
	id := context.GalleryID(r.Context())
	if id == nil {
		WriteErrJSON(w, http.StatusBadRequest, "Invalid ID. Must be an integer.")
		return
	}

	res, err := c.GalleryService.ByID(galleryM.GetGalleryByIDRequest{ID: *id})
	if err != nil {
		if errors.Is(err, galleryM.ErrGalleryNotFound) {
			WriteErrJSON(w, http.StatusNotFound, fmt.Sprintf("Gallery with id = %d not found", *id))
			return
		}

		WriteErrJSON(w, http.StatusInternalServerError, "Failed to get gallery")
		return
	}

	WriteJSON(w, http.StatusOK, res.Gallery.ToJSON())
}

func (c *Galleries) SetGalleryTitle(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Title string `json:"title"`
		ID    int    `json:"id"`
	}
	var in input
	if err := ReadJSON(r, &in); err != nil {
		WriteErrJSON(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}

	id := context.GalleryID(r.Context())
	if id == nil {
		WriteErrJSON(w, http.StatusBadRequest, "Invalid ID. Must be an integer.")
		return
	}
	in.ID = *id
	res, err := c.GalleryService.SetTitle(galleryM.SetGalleryTitleRequest{ID: in.ID, Title: in.Title})
	if err != nil {
		if errors.Is(err, galleryM.ErrGalleryNotFound) {
			WriteErrJSON(w, http.StatusNotFound, fmt.Sprintf("Gallery with id = %d not found", in.ID))
			return
		}

		WriteErrJSON(w, http.StatusInternalServerError, "Failed to set gallery title")
		return
	}

	WriteJSON(w, http.StatusOK, res.Gallery.ToJSON())
}

func (c *Galleries) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	id := context.GalleryID(r.Context())
	if id == nil {
		WriteErrJSON(w, http.StatusBadRequest, "Invalid ID. Must be an integer.")
		return
	}

	err := c.GalleryService.DeleteByID(galleryM.DeleteGalleryRequest{ID: *id})
	if err != nil {
		fmt.Println(err)
		WriteErrJSON(w, http.StatusInternalServerError, "Failed to delete gallery")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Gallery deleted"})
}
