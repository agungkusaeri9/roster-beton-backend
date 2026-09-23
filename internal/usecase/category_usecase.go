package usecase

import (
	"go-arch/internal/entity"
	"go-arch/internal/repository"
)

type CategoryUsecase interface {
	GetAll() ([]*entity.Category, error)
	GetBySlug(slug string) (*entity.Category, error)
	GetByID(id int64) (*entity.Category, error)
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUsecase) GetAll() ([]*entity.Category, error) {
	return u.categoryRepo.FindAll()
}

func (u *categoryUsecase) GetBySlug(slug string) (*entity.Category, error) {
	return u.categoryRepo.FindBySlug(slug)
}

func (u *categoryUsecase) GetByID(id int64) (*entity.Category, error) {
	return u.categoryRepo.FindByID(id)
}
