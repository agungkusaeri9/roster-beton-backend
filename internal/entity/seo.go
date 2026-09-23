package entity

import "time"

// SEO represents polymorphic SEO metadata for products, articles, categories, pages, etc.
type SEO struct {
	ID             int64     `db:"id" json:"id"`
	ReferenceType  string    `db:"reference_type" json:"reference_type"` // e.g. "product", "article", "category", "page"
	ReferenceID    int64     `db:"reference_id" json:"reference_id"`     // ID dari entitas terkait
	MetaTitle      string    `db:"meta_title" json:"meta_title,omitempty"`
	MetaDescription string   `db:"meta_description" json:"meta_description,omitempty"`
	MetaKeywords   string    `db:"meta_keywords" json:"meta_keywords,omitempty"`
	CanonicalURL   string    `db:"canonical_url" json:"canonical_url,omitempty"`
	OGTitle        string    `db:"og_title" json:"og_title,omitempty"`
	OGDescription  string    `db:"og_description" json:"og_description,omitempty"`
	OGImage        string    `db:"og_image" json:"og_image,omitempty"`
	OGType         string    `db:"og_type" json:"og_type,omitempty"`
	StructuredData *string   `db:"structured_data" json:"structured_data,omitempty"` // Schema.org JSON-LD
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
