package repository

import "go-arch/internal/entity"

type CategoryRepository interface {
	FindAll() ([]*entity.Category, error)
	FindBySlug(slug string) (*entity.Category, error)
	FindByID(id int64) (*entity.Category, error)
}
