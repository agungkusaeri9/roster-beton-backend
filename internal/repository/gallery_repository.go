package repository

import "go-arch/internal/entity"

type GalleryRepository interface {
	FindAll(category string) ([]*entity.Gallery, error)
	FindByID(id int64) (*entity.Gallery, error)
}
