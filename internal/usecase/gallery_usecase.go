package usecase

import (
	"go-arch/internal/entity"
	"go-arch/internal/repository"
)

type GalleryUsecase interface {
	GetAll(category string) ([]*entity.Gallery, error)
	GetByID(id int64) (*entity.Gallery, error)
}

type galleryUsecase struct {
	repo repository.GalleryRepository
}

func NewGalleryUsecase(repo repository.GalleryRepository) GalleryUsecase {
	return &galleryUsecase{repo: repo}
}

func (uc *galleryUsecase) GetAll(category string) ([]*entity.Gallery, error) {
	return uc.repo.FindAll(category)
}

func (uc *galleryUsecase) GetByID(id int64) (*entity.Gallery, error) {
	return uc.repo.FindByID(id)
}
