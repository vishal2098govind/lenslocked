package context

import "context"

var (
	galleryIdKey key = "galleryId"
)

func WithGalleryID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, galleryIdKey, id)
}

func GalleryID(ctx context.Context) *int {
	id, ok := ctx.Value(galleryIdKey).(int)
	if !ok {
		return nil
	}
	return &id
}
