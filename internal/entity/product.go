package entity

import "time"

// Product represents the main product entity (Clean, tanpa kolom SEO langsung)
type Product struct {
	ID          int64            `db:"id" json:"id"`
	Name        string           `db:"name" json:"name"`
	Slug        string           `db:"slug" json:"slug"`
	CategoryID  *int64           `db:"category_id" json:"category_id,omitempty"`
	Category    *Category        `db:"-" json:"category,omitempty"` // Relasi ke entity Category
	Description string           `db:"description" json:"description"`
	Image       string           `db:"image" json:"image"` // Main thumbnail URL
	Price       *float64         `db:"price" json:"price,omitempty"`
	Dimension   *string          `db:"dimension" json:"dimension,omitempty"` // e.g. "20x20x10 cm"
	Weight      *float64         `db:"weight" json:"weight,omitempty"`
	Stock       int              `db:"stock" json:"stock,omitempty"`
	IsActive    bool             `db:"is_active" json:"is_active"`
	IsFeatured  bool             `db:"is_featured" json:"is_featured"`

	// Relations
	Galleries []ProductGallery `db:"-" json:"galleries,omitempty"`
	Gallery   []string         `db:"-" json:"gallery,omitempty"` // Array string format frontend
	SEO       *SEO             `db:"-" json:"seo,omitempty"`     // Polymorphic SEO relation

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
