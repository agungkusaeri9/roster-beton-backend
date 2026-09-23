package pgsql

import (
	"fmt"

	"go-arch/internal/entity"
	"go-arch/internal/repository"

	"github.com/jmoiron/sqlx"
)

type categoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepoPg(db *sqlx.DB) repository.CategoryRepository {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) FindAll() ([]*entity.Category, error) {
	query := `
		SELECT id, name, slug, description, created_at, updated_at
		FROM categories
		ORDER BY id ASC
	`
	var categories []*entity.Category
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	if categories == nil {
		categories = []*entity.Category{}
	}

	return categories, nil
}

func (r *categoryRepo) FindBySlug(slug string) (*entity.Category, error) {
	query := `
		SELECT id, name, slug, description, created_at, updated_at
		FROM categories
		WHERE slug = $1
		LIMIT 1
	`
	var category entity.Category
	err := r.db.Get(&category, query, slug)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepo) FindByID(id int64) (*entity.Category, error) {
	query := `
		SELECT id, name, slug, description, created_at, updated_at
		FROM categories
		WHERE id = $1
		LIMIT 1
	`
	var category entity.Category
	err := r.db.Get(&category, query, id)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
