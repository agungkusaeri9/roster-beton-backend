package pgsql

import (
	"database/sql"
	"fmt"
	"go-arch/internal/entity"
	"go-arch/internal/repository"
	"strings"

	"github.com/jmoiron/sqlx"
)

type galleryRepo struct {
	db *sqlx.DB
}

func NewGalleryRepoPg(db *sqlx.DB) repository.GalleryRepository {
	return &galleryRepo{db: db}
}

func (r *galleryRepo) FindAll(category string) ([]*entity.Gallery, error) {
	query := `
		SELECT id, src, title, category, COALESCE(alt_text, '') as alt_text, sort_order, created_at, updated_at
		FROM galleries
	`
	var args []interface{}
	if strings.TrimSpace(category) != "" && strings.ToLower(category) != "semua" {
		query += " WHERE LOWER(category) = LOWER($1)"
		args = append(args, strings.TrimSpace(category))
	}
	query += " ORDER BY sort_order ASC, id ASC"

	var galleries []*entity.Gallery
	err := r.db.Select(&galleries, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed fetching galleries: %w", err)
	}

	if galleries == nil {
		galleries = []*entity.Gallery{}
	}
	return galleries, nil
}

func (r *galleryRepo) FindByID(id int64) (*entity.Gallery, error) {
	query := `
		SELECT id, src, title, category, COALESCE(alt_text, '') as alt_text, sort_order, created_at, updated_at
		FROM galleries
		WHERE id = $1
	`
	var item entity.Gallery
	err := r.db.Get(&item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed fetching gallery by id: %w", err)
	}
	return &item, nil
}
