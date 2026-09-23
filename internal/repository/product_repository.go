package repository

import "go-arch/internal/entity"

type ProductFilter struct {
	Page     int
	Limit    int
	Category string
	Search   string
}

type ProductRepository interface {
	FindAll(filter ProductFilter) ([]*entity.Product, int64, error)
	FindBySlug(slug string) (*entity.Product, error)
	FindByID(id int64) (*entity.Product, error)
}
