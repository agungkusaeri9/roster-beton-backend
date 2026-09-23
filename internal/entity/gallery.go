package entity

import "time"

// Gallery represents general gallery items (projects, factory, products, company) matching frontend data/gallery.json
type Gallery struct {
	ID        int64     `db:"id" json:"id"`
	Src       string    `db:"src" json:"src"`
	Title     string    `db:"title" json:"title"`
	Category  string    `db:"category" json:"category"` // "Produk", "Proyek", "Produksi", "Layanan", "Perusahaan"
	AltText   string    `db:"alt_text" json:"alt_text,omitempty"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
