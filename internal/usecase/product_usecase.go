package usecase

import (
	"math"

	"go-arch/internal/entity"
	"go-arch/internal/repository"
)

type ProductUsecase interface {
	GetAll(page, limit int, category, search string) ([]*entity.Product, entity.PaginationMeta, error)
	GetBySlug(slug string) (*entity.Product, error)
	GetByID(id int64) (*entity.Product, error)
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) GetAll(page, limit int, category, search string) ([]*entity.Product, entity.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	filter := repository.ProductFilter{
		Page:     page,
		Limit:    limit,
		Category: category,
		Search:   search,
	}

	products, total, err := u.productRepo.FindAll(filter)
	if err != nil {
		return nil, entity.PaginationMeta{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	pagination := entity.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalData:   total,
		TotalPages:  totalPages,
	}

	return products, pagination, nil
}

func (u *productUsecase) GetBySlug(slug string) (*entity.Product, error) {
	return u.productRepo.FindBySlug(slug)
}

func (u *productUsecase) GetByID(id int64) (*entity.Product, error) {
	return u.productRepo.FindByID(id)
}
