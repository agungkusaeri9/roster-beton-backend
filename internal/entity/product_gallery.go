package entity

import "time"

// ProductGallery represents additional images associated with a product
type ProductGallery struct {
	ID        int64     `db:"id" json:"id"`
	ProductID int64     `db:"product_id" json:"product_id"`
	ImageURL  string    `db:"image_url" json:"image_url"`
	AltText   string    `db:"alt_text" json:"alt_text,omitempty"` // Image SEO & Accessibility
	SortOrder int       `db:"sort_order" json:"sort_order"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
