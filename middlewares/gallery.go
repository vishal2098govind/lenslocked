package middlewares

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vishal2098govind/lenslocked/context"
	galleryM "github.com/vishal2098govind/lenslocked/models/gallery"
)

type GalleryMW struct {
	GalleryService *galleryM.GalleryService
}

func (gmw *GalleryMW) SetGalleryID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := chi.URLParam(r, "id")
		id, err := strconv.Atoi(ids)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithGalleryID(ctx, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
